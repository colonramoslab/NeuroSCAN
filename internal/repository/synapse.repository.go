package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"neuroscan/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SynapseRepository interface {
	GetSynapseByUID(ctx context.Context, uid string, timepoint int) (domain.Synapse, error)
	SynapseExists(ctx context.Context, uid string, timepoint int) (bool, error)
	SearchSynapses(ctx context.Context, query domain.APIV1Request) ([]domain.Synapse, error)
	CountSynapses(ctx context.Context, query domain.APIV1Request) (int, error)
	CreateSynapse(ctx context.Context, synapse domain.Synapse) error
	DeleteSynapse(ctx context.Context, uid string, timepoint int) error
	IngestSynapse(ctx context.Context, synapse domain.Synapse, skipExisting bool, force bool) (bool, error)
	TruncateSynapses(ctx context.Context) error
}

type PostgresSynapseRepository struct {
	DB *pgxpool.Pool
}

func NewPostgresSynapseRepository(db *pgxpool.Pool) *PostgresSynapseRepository {
	return &PostgresSynapseRepository{
		DB: db,
	}
}

func (r *PostgresSynapseRepository) GetSynapseByUID(ctx context.Context, uid string, timepoint int) (domain.Synapse, error) {
	query := "SELECT id, uid, ulid, timepoint, synapse_type, filename, color FROM synapses WHERE uid = $1 AND timepoint = $2"

	var synapse domain.Synapse
	err := r.DB.QueryRow(ctx, query, uid, timepoint).Scan(&synapse.ID, &synapse.UID, &synapse.ULID, &synapse.Timepoint, &synapse.SynapseType, &synapse.Filename, &synapse.Color)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Synapse{}, nil
		}

		return domain.Synapse{}, err
	}

	return synapse, nil
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

	synapses, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Synapse])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.Synapse{}, nil
		}

		return nil, err
	}

	return synapses, nil
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

func (r *PostgresSynapseRepository) ParseSynapseAPIV1Request(ctx context.Context, req domain.APIV1Request) (string, []interface{}) {

	queryParts := []string{"where 1=1"}
	args := []interface{}{}

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

		for _, synapseType := range synapseTypes {
			if synapseType == "chemical" {
				containsChemical = true
				break
			}
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
