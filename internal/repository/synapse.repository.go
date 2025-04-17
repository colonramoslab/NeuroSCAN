package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	"neuroscan/internal/cache"
	"neuroscan/internal/domain"
	"neuroscan/internal/toolshed"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SynapseRepository interface {
	GetSynapseByULID(ctx context.Context, id string) (domain.Synapse, error)
	GetSynapseByUID(ctx context.Context, uid string, timepoint int) (domain.Synapse, error)
	SynapseExists(ctx context.Context, uid string, timepoint int) (bool, error)
	SearchSynapses(ctx context.Context, query domain.APIV1Request) ([]domain.Synapse, error)
	CountSynapses(ctx context.Context, query domain.APIV1Request) (int, error)
	CreateSynapse(ctx context.Context, synapse domain.Synapse) error
	DeleteSynapse(ctx context.Context, uid string, timepoint int) error
	IngestSynapse(ctx context.Context, synapse domain.Synapse, skipExisting bool, force bool) (bool, error)
	TruncateSynapses(ctx context.Context) error
}

type Synapse struct {
	ID          int                `db:"id"`
	ULID        string             `db:"ulid"`
	UID         string             `db:"uid"`
	Timepoint   int                `db:"timepoint"`
	SynapseType domain.SynapseType `db:"type"`
	Filename    string             `db:"filename"`
	Color       toolshed.Color     `db:"color"`
}

func (s *Synapse) ToDomain(tnrs *int, synapses *[]domain.SynapseItem) domain.Synapse {
	synapse := domain.Synapse{
		ID:          s.ID,
		ULID:        s.ULID,
		UID:         s.UID,
		Timepoint:   s.Timepoint,
		SynapseType: s.SynapseType,
		Filename:    s.Filename,
		Color:       s.Color,
	}

	if tnrs != nil {
		synapse.TotalNRSynapses = tnrs
	}

	if synapses != nil {
		synapse.Synapses = synapses
	}

	return synapse
}

type PostgresSynapseRepository struct {
	cache cache.Cache
	DB    *pgxpool.Pool
}

func NewPostgresSynapseRepository(db *pgxpool.Pool, c cache.Cache) *PostgresSynapseRepository {
	return &PostgresSynapseRepository{
		cache: c,
		DB:    db,
	}
}

func (r *PostgresSynapseRepository) GetSynapseByULID(ctx context.Context, id string) (domain.Synapse, error) {
	query := "SELECT * FROM synapses WHERE ulid = $1"

	var synapse Synapse
	err := r.DB.QueryRow(ctx, query, id).Scan(&synapse.ID, &synapse.UID, &synapse.ULID, &synapse.Timepoint, &synapse.SynapseType, &synapse.Filename, &synapse.Color)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Synapse{}, nil
		}

		return domain.Synapse{}, err
	}

	tnrs, err := r.NerveRingSynapseCount(ctx, synapse.Timepoint)
	if err != nil {
		return domain.Synapse{}, err
	}

	return synapse.ToDomain(&tnrs), nil
}

func (r *PostgresSynapseRepository) GetSynapseByUID(ctx context.Context, uid string, timepoint int) (domain.Synapse, error) {
	query := "SELECT id, uid, ulid, timepoint, synapse_type, filename, color FROM synapses WHERE uid = $1 AND timepoint = $2"

	var synapse Synapse
	err := r.DB.QueryRow(ctx, query, uid, timepoint).Scan(&synapse.ID, &synapse.UID, &synapse.ULID, &synapse.Timepoint, &synapse.SynapseType, &synapse.Filename, &synapse.Color)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Synapse{}, nil
		}

		return domain.Synapse{}, err
	}

	return synapse.ToDomain(), nil
}

func (r *PostgresSynapseRepository) SynapseExists(ctx context.Context, uid string, timepoint int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM synapses WHERE uid = $1 AND timepoint = $2)"

	var exists bool
	err := r.DB.QueryRow(ctx, query, uid, timepoint).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return exists, nil
}

func (r *PostgresSynapseRepository) SearchSynapses(ctx context.Context, query domain.APIV1Request) ([]domain.Synapse, error) {
	q := "SELECT id, uid, ulid, timepoint, synapse_type, filename, color FROM synapses "

	parsedQuery, args := r.ParseSynapseAPIV1Request(ctx, query)

	q += parsedQuery

	rows, _ := r.DB.Query(ctx, q, args...)

	synapses, err := pgx.CollectRows(rows, pgx.RowToStructByName[Synapse])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.Synapse{}, nil
		}

		return nil, err
	}

	domainSynapses := make([]domain.Synapse, len(synapses))

	for i := range synapses {
		domainSynapses[i] = synapses[i].ToDomain(nil, nil)
	}

	return domainSynapses, nil
}

func (r *PostgresSynapseRepository) CountSynapses(ctx context.Context, query domain.APIV1Request) (int, error) {
	var count int

	q := "SELECT COUNT(*) FROM synapses "

	parsedQuery, args := r.ParseSynapseAPIV1Request(ctx, query)

	q += parsedQuery

	err := r.DB.QueryRow(ctx, q, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresSynapseRepository) SynapseCount(ctx context.Context, uid string, timepoint int) ([]domain.SynapseItem, error) {
	parts := strings.Split(uid, "&")
	prefix := parts[0]
	like := fmt.Sprintf("%s%%", prefix)
	cacheKey := fmt.Sprintf("synapse:synapse_count:%s:%d", prefix, timepoint)

	if cachedSynapseCount, found := r.cache.Get(cacheKey); found {
		if cached, ok := cachedSynapseCount.([]domain.SynapseItem); ok {
			return cached, nil
		}
	}

	query := "SELECT split_part(uid, '~', 1) AS syn_identity, COUNT(*) AS total FROM synapses WHERE uid LIKE $1 AND timepoint = $2 GROUP BY syn_identity ORDER BY syn_identity ASC;"

	rows, err := r.DB.Query(ctx, query, like, timepoint)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.SynapseItem{}, nil
		}

		return []domain.SynapseItem{}, err
	}

	defer rows.Close()

	synapses := []domain.SynapseItem{}
	for rows.Next() {
		var synapse domain.SynapseItem
		err := rows.Scan(&synapse.Name, &synapse.Count)
		if err != nil {
			return []domain.SynapseItem{}, err
		}
		synapses = append(synapses, synapse)
	}

	r.cache.Set(cacheKey, synapses)

	return synapses, nil
}

func (r *PostgresSynapseRepository) NerveRingSynapseCount(ctx context.Context, timepoint int) (int, error) {
	cacheKey := fmt.Sprintf("nervering:total_synapses:%d", timepoint)

	if cachedTNRS, found := r.cache.Get(cacheKey); found {
		if cached, ok := cachedTNRS.(int); ok {
			return cached, nil
		}
	}

	query := "SELECT count(*) FROM synapses WHERE timepoint = $1;"

	var total sql.NullInt64
	err := r.DB.QueryRow(ctx, query, timepoint).Scan(&total)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	if total.Valid {
		count := total.Int64

		r.cache.Set(cacheKey, count)

		return int(count), nil
	}

	return 0, nil
}

func (r *PostgresSynapseRepository) CreateSynapse(ctx context.Context, synapse domain.Synapse) error {
	exists, err := r.SynapseExists(ctx, synapse.UID, synapse.Timepoint)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("synapse already exists")
	}

	query := "INSERT INTO synapses (uid, ulid, timepoint, synapse_type, filename, color) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING"

	_, err = r.DB.Exec(ctx, query, synapse.UID, synapse.ULID, synapse.Timepoint, synapse.SynapseType, synapse.Filename, synapse.Color)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresSynapseRepository) DeleteSynapse(ctx context.Context, uid string, timepoint int) error {
	query := "DELETE FROM synapses WHERE uid = $1 AND timepoint = $2"

	_, err := r.DB.Exec(ctx, query, uid, timepoint)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresSynapseRepository) IngestSynapse(ctx context.Context, synapse domain.Synapse, skipExisting bool, force bool) (bool, error) {
	exists, err := r.SynapseExists(ctx, synapse.UID, synapse.Timepoint)
	if err != nil {
		return false, err
	}

	if skipExisting && exists {
		return true, nil
	}

	if force && exists {
		err = r.DeleteSynapse(ctx, synapse.UID, synapse.Timepoint)
		if err != nil {
			return false, err
		}
	}

	err = r.CreateSynapse(ctx, synapse)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *PostgresSynapseRepository) TruncateSynapses(ctx context.Context) error {
	query := "TRUNCATE TABLE synapses RESTART IDENTITY CASCADE"

	_, err := r.DB.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresSynapseRepository) ParseSynapseAPIV1Request(ctx context.Context, req domain.APIV1Request) (string, []any) {
	queryParts := []string{"where 1=1"}
	args := []any{}

	if req.Timepoint != nil {
		args = append(args, req.Timepoint)
		queryParts = append(queryParts, fmt.Sprintf("timepoint = $%d", len(args)))
	}

	if len(req.UIDs) > 0 {
		// we need to build a query where UID is like or, looping over the UIDs, wrapping them in % and adding them to the array[]
		uidArray := []string{}
		for _, uid := range req.UIDs {
			uidArray = append(uidArray, fmt.Sprintf("%%%s%%", strings.ToLower(uid)))
		}
		args = append(args, uidArray)
		queryParts = append(queryParts, fmt.Sprintf("LOWER(uid) ILIKE ANY($%d)", len(args)))
	}

	if len(req.Types) > 0 {
		synapseTypes := req.Types
		containsChemical := false
		args = append(args, req.Types)

		if slices.Contains(synapseTypes, "chemical") {
			containsChemical = true
		}

		if containsChemical {
			queryParts = append(queryParts, fmt.Sprintf("synapse_type = ANY($%d) OR synapse_type IS NULL", len(args)))
		} else {
			queryParts = append(queryParts, fmt.Sprintf("synapse_type = ANY($%d)", len(args)))
		}
	}

	if req.PreNeuron != "" {
		args = append(args, strings.ToLower(req.PreNeuron))
		queryParts = append(queryParts, fmt.Sprintf("AND LOWER(uid) LIKE '$%d%%'", len(args)))
	}

	if req.PostNeuron != "" {
		args = append(args, strings.Join(req.Types, "|"))
		args = append(args, strings.ToLower(req.PostNeuron))
		queryParts = append(queryParts, fmt.Sprintf("AND LOWER(uid) SIMILAR TO '%%($%d%%$%d)%%|~%%)'", len(args)-1, len(args)))
	}

	query := strings.Join(queryParts, " AND ")

	// if count is true, return the query and args before adding the sort and limit
	if req.Count {
		return query, args
	}

	if req.Sort != "" {
		// split by ":", first part is the field, second part is the direction
		parts := strings.Split(req.Sort, ":")

		if len(parts) == 2 {

			// if the second part is not asc or desc, default to asc
			if parts[1] != "asc" && parts[1] != "desc" {
				parts[1] = "asc"
			}

			query += fmt.Sprintf(" order by %s %s", parts[0], parts[1])
		}
	}

	if req.Limit > 0 {
		args = append(args, req.Limit)
		query += fmt.Sprintf(" limit $%d", len(args))
	} else {
		query += " limit 100"
	}

	if req.Offset > 0 {
		args = append(args, req.Offset)
		query += fmt.Sprintf(" offset $%d", len(args))
	}

	return query, args
}
