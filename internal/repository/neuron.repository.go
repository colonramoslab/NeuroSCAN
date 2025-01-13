package repository

import (
	"context"
	"errors"
	"fmt"

	"neuroscan/internal/domain"
	"neuroscan/internal/toolshed"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NeuronQuery struct {
	IDs			 []string
	UIDs		 []string
	Timepoint	 *int
	Limit 		 *int
	Offset 		 *int
}

type NeuronRepository interface {
	GetAllNeurons(ctx context.Context, timepoint int) ([]domain.Neuron, error)
	GetNeuronByID(ctx context.Context, uid string, timepoint int) (domain.Neuron, error)
	GetNeuronByUID(ctx context.Context, uid string, timepoint int) (domain.Neuron, error)
	NeuronExists(ctx context.Context, uid string, timepoint int) (bool, error)
	SearchNeurons(ctx context.Context, query NeuronQuery) ([]domain.Neuron, error)
	CreateNeuron(ctx context.Context, uid string, filename string, timepoint int, color toolshed.Color) error
	DeleteNeuron(ctx context.Context, uid string, timepoint int) error
	IngestNeuron(ctx context.Context, neuron domain.Neuron, skipExisting bool, force bool) (bool, error)
}

type PostgresNeuronRepository struct {
	DB *pgxpool.Pool
}

func NewPostgresNeuronRepository(db *pgxpool.Pool) *PostgresNeuronRepository {
	return &PostgresNeuronRepository{
		DB: db,
	}
}

func (r *PostgresNeuronRepository) GetAllNeurons(ctx context.Context, timepoint int) ([]domain.Neuron, error) {
	query := "SELECT id, uid, timepoint, filename, color FROM neurons WHERE timepoint = $1"

	rows, _ := r.DB.Query(ctx, query, timepoint)

	neurons, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Neuron])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.Neuron{}, nil
		}

		return nil, err
	}

	return neurons, nil
}

func (r *PostgresNeuronRepository) GetNeuronByID(ctx context.Context, uid string, timepoint int) (domain.Neuron, error) {
	query := "SELECT id, uid, timepoint, filename, color FROM neurons WHERE id = $1 AND timepoint = $2"

	var neuron domain.Neuron
	err := r.DB.QueryRow(ctx, query, uid, timepoint).Scan(&neuron.ID, &neuron.UID, &neuron.Timepoint, &neuron.Filename, &neuron.Color)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Neuron{}, nil
		}

		return domain.Neuron{}, err
	}

	return neuron, nil
}

func (r *PostgresNeuronRepository) GetNeuronByUID(ctx context.Context, uid string, timepoint int) (domain.Neuron, error) {
	query := "SELECT id, uid, timepoint, filename, color FROM neurons WHERE uid = $1 AND timepoint = $2"

	var neuron domain.Neuron
	err := r.DB.QueryRow(ctx, query, uid, timepoint).Scan(&neuron.ID, &neuron.UID, &neuron.Timepoint, &neuron.Filename, &neuron.Color)
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

func (r *PostgresNeuronRepository) SearchNeurons(ctx context.Context, query NeuronQuery) ([]domain.Neuron, error) {
	q := "SELECT id, uid, timepoint, filename, color FROM neurons WHERE 1=1"
	args := []interface{}{}

	if query.Timepoint != nil {
		q += fmt.Sprintf(" AND timepoint = $%d", len(args)+1)
		args = append(args, query.Timepoint)
	}

	if len(query.IDs) > 0 {
		q += fmt.Sprintf(" AND id = ANY($%d)", len(args)+1)
		args = append(args, query.IDs)
	}

	if len(query.UIDs) > 0 {
		q += fmt.Sprintf(" AND uid = ANY($%d)", len(args)+1)
		args = append(args, query.UIDs)
	}

	if query.Limit != nil {
		q += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, query.Limit)
	}

	if query.Offset != nil {
		q += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, query.Offset)
	}

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

func (r *PostgresNeuronRepository) CreateNeuron(ctx context.Context, uid string, filename string, timepoint int, color toolshed.Color) error {
	exists, err := r.NeuronExists(ctx, uid, timepoint)

	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("neuron already exists")
	}

	query := "INSERT INTO neurons (uid, timepoint, filename, color) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING"

	_, err = r.DB.Exec(ctx, query, uid, timepoint, filename, color)
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

	err = r.CreateNeuron(ctx, neuron.UID, neuron.Filename, neuron.Timepoint, neuron.Color)
	if err != nil {
		return false, err
	}

	return true, nil
}
