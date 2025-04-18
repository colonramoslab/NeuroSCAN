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

type NeuronRepository interface {
	GetNeuronByULID(ctx context.Context, id string) (domain.Neuron, error)
	GetNeuronByUID(ctx context.Context, uid string, timepoint int) (domain.Neuron, error)
	NeuronExists(ctx context.Context, uid string, timepoint int) (bool, error)
	SearchNeurons(ctx context.Context, query domain.APIV1Request) ([]domain.Neuron, error)
	CountNeurons(ctx context.Context, query domain.APIV1Request) (int, error)
	CreateNeuron(ctx context.Context, neuron domain.Neuron) error
	DeleteNeuron(ctx context.Context, uid string, timepoint int) error
	IngestNeuron(ctx context.Context, neuron domain.Neuron, skipExisting bool, force bool) (bool, error)
	TruncateNeurons(ctx context.Context) error
}

type PostgresNeuronRepository struct {
	DB *pgxpool.Pool
}

func NewPostgresNeuronRepository(db *pgxpool.Pool) *PostgresNeuronRepository {
	return &PostgresNeuronRepository{
		DB: db,
	}
}

func (r *PostgresNeuronRepository) GetNeuronByULID(ctx context.Context, id string) (domain.Neuron, error) {
	query := "SELECT * FROM neurons WHERE ulid = $1"

	var neuron domain.Neuron
	err := r.DB.QueryRow(ctx, query, id).Scan(&neuron.ID, &neuron.UID, &neuron.ULID, &neuron.Timepoint, &neuron.Filename, &neuron.Color)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Neuron{}, nil
		}

		return domain.Neuron{}, err
	}

	return neuron, nil
}

func (r *PostgresNeuronRepository) GetNeuronByUID(ctx context.Context, uid string, timepoint int) (domain.Neuron, error) {
	query := "SELECT id, uid, ulid, timepoint, filename, color FROM neurons WHERE uid = $1 AND timepoint = $2"

	var neuron domain.Neuron
	err := r.DB.QueryRow(ctx, query, uid, timepoint).Scan(&neuron.ID, &neuron.UID, &neuron.ULID, &neuron.Timepoint, &neuron.Filename, &neuron.Color)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Neuron{}, nil
		}

		return domain.Neuron{}, err
	}

	return neuron, nil
}

func (r *PostgresNeuronRepository) NeuronExists(ctx context.Context, uid string, timepoint int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM neurons WHERE uid = $1 AND timepoint = $2)"

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

func (r *PostgresNeuronRepository) SearchNeurons(ctx context.Context, query domain.APIV1Request) ([]domain.Neuron, error) {
	q := "SELECT id, uid, ulid, timepoint, filename, color FROM neurons "

	parsedQuery, args := r.ParseNeuronAPIV1Request(ctx, query)

	q += parsedQuery

	rows, _ := r.DB.Query(ctx, q, args...)

	neurons, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Neuron])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.Neuron{}, nil
		}

		return nil, err
	}

	return neurons, err
}

func (r *PostgresNeuronRepository) CountNeurons(ctx context.Context, query domain.APIV1Request) (int, error) {
	var count int

	q := "SELECT COUNT(*) FROM neurons "

	parsedQuery, args := r.ParseNeuronAPIV1Request(ctx, query)

	q += parsedQuery

	err := r.DB.QueryRow(ctx, q, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresNeuronRepository) CreateNeuron(ctx context.Context, neuron domain.Neuron) error {
	exists, err := r.NeuronExists(ctx, neuron.UID, neuron.Timepoint)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("neuron already exists")
	}

	query := "INSERT INTO neurons (uid, ulid, timepoint, filename, color) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING"

	_, err = r.DB.Exec(ctx, query, neuron.UID, neuron.ULID, neuron.Timepoint, neuron.Filename, neuron.Color)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresNeuronRepository) DeleteNeuron(ctx context.Context, uid string, timepoint int) error {
	query := "DELETE FROM neurons WHERE uid = $1 AND timepoint = $2"

	_, err := r.DB.Exec(ctx, query, uid, timepoint)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresNeuronRepository) IngestNeuron(ctx context.Context, neuron domain.Neuron, skipExisting bool, force bool) (bool, error) {
	exists, err := r.NeuronExists(ctx, neuron.UID, neuron.Timepoint)
	if err != nil {
		return false, err
	}

	if skipExisting && exists {
		return true, nil
	}

	if force && exists {
		err = r.DeleteNeuron(ctx, neuron.UID, neuron.Timepoint)
		if err != nil {
			return false, err
		}
	}

	err = r.CreateNeuron(ctx, neuron)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *PostgresNeuronRepository) TruncateNeurons(ctx context.Context) error {
	query := "TRUNCATE TABLE neurons RESTART IDENTITY CASCADE"

	_, err := r.DB.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresNeuronRepository) ParseNeuronAPIV1Request(ctx context.Context, req domain.APIV1Request) (string, []interface{}) {
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
