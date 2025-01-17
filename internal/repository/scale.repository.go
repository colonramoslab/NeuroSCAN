package repository

import (
	"context"
	"errors"
	"fmt"

	"neuroscan/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScaleRepository interface {
	GetScaleByTimepoint(ctx context.Context, timepoint int) (domain.Scale, error)
	ScaleExists(ctx context.Context, timepoint int) (bool, error)
	CreateScale(ctx context.Context, scale domain.Scale) error
	DeleteScale(ctx context.Context, timepoint int) error
	IngestScale(ctx context.Context, scale domain.Scale, skipExisting bool, force bool) (bool, error)
	TruncateScales(ctx context.Context) error
}

type PostgresScaleRepository struct {
	DB *pgxpool.Pool
}

func NewPostgresScaleRepository(db *pgxpool.Pool) *PostgresScaleRepository {
	return &PostgresScaleRepository{
		DB: db,
	}
}

func (r *PostgresScaleRepository) GetScaleByTimepoint(ctx context.Context, timepoint int) (domain.Scale, error) {
	query := "SELECT id, uid, timepoint, filename, color FROM scales WHERE timepoint = $1"

	var scale domain.Scale
	err := r.DB.QueryRow(ctx, query, timepoint).Scan(&scale.ID, &scale.UID, &scale.Timepoint, &scale.Filename, &scale.Color)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Scale{}, nil
		}

		return domain.Scale{}, err
	}

	return scale, nil
}

func (r *PostgresScaleRepository) ScaleExists(ctx context.Context, timepoint int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM scales WHERE timepoint = $1)"

	var exists bool
	err := r.DB.QueryRow(ctx, query, timepoint).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return exists, nil
}

func (r *PostgresScaleRepository) CreateScale(ctx context.Context, scale domain.Scale) error {
	exists, err := r.ScaleExists(ctx, scale.Timepoint)

	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("scale already exists")
	}

	query := "INSERT INTO scales (uid, timepoint, filename, color) VALUES ($1, $2, $3, $4)"

	_, err = r.DB.Exec(ctx, query, scale.UID, scale.Timepoint, scale.Filename, scale.Color)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresScaleRepository) DeleteScale(ctx context.Context, timepoint int) error {
	query := "DELETE FROM scales WHERE timepoint = $1"

	_, err := r.DB.Exec(ctx, query, timepoint)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresScaleRepository) IngestScale(ctx context.Context, scale domain.Scale, skipExisting bool, force bool) (bool, error) {
	exists, err := r.ScaleExists(ctx, scale.Timepoint)

	if err != nil {
		return false, err
	}

	if skipExisting && exists {
		return true, nil
	}

	if force && exists {
		err = r.DeleteScale(ctx, scale.Timepoint)
		if err != nil {
			return false, err
		}
	}

	err = r.CreateScale(ctx, scale)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *PostgresScaleRepository) TruncateScales(ctx context.Context) error {
	query := "TRUNCATE TABLE scales RESTART IDENTITY CASCADE"

	_, err := r.DB.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
