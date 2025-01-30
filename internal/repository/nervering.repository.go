package repository

import (
	"context"
	"errors"
	"fmt"

	"neuroscan/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NerveRingRepository interface {
	GetNerveRingByTimepoint(ctx context.Context, timepoint int) (domain.NerveRing, error)
	NerveRingExists(ctx context.Context, timepoint int) (bool, error)
	CreateNerveRing(ctx context.Context, nerveRing domain.NerveRing) error
	DeleteNerveRing(ctx context.Context, uid string, timepoint int) error
	IngestNerveRing(ctx context.Context, nerveRing domain.NerveRing, skipExisting bool, force bool) (bool, error)
	TruncateNerveRings(ctx context.Context) error
}

type PostgresNerveRingRepository struct {
	DB *pgxpool.Pool
}

func NewPostgresNerveRingRepository(db *pgxpool.Pool) *PostgresNerveRingRepository {
	return &PostgresNerveRingRepository{
		DB: db,
	}
}

func (r *PostgresNerveRingRepository) GetNerveRingByTimepoint(ctx context.Context, timepoint int) (domain.NerveRing, error) {
	query := "SELECT id, uid, ulid, timepoint, filename, color FROM nerve_rings WHERE timepoint = $1"

	var nerveRing domain.NerveRing
	err := r.DB.QueryRow(ctx, query, timepoint).Scan(&nerveRing.ID, &nerveRing.UID, &nerveRing.ULID, &nerveRing.Timepoint, &nerveRing.Filename, &nerveRing.Color)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.NerveRing{}, nil
		}

		return domain.NerveRing{}, err
	}

	return nerveRing, nil
}

func (r *PostgresNerveRingRepository) NerveRingExists(ctx context.Context, timepoint int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM nerve_rings WHERE timepoint = $1)"

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

func (r *PostgresNerveRingRepository) CreateNerveRing(ctx context.Context, nerveRing domain.NerveRing) error {
	exists, err := r.NerveRingExists(ctx, nerveRing.Timepoint)

	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("nerve ring already exists")
	}

	query := "INSERT INTO nerve_rings (uid, ulid, timepoint, filename, color) VALUES ($1, $2, $3, $4, $5)"

	_, err = r.DB.Exec(ctx, query, nerveRing.UID, nerveRing.ULID, nerveRing.Timepoint, nerveRing.Filename, nerveRing.Color)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresNerveRingRepository) DeleteNerveRing(ctx context.Context, uid string, timepoint int) error {
	query := "DELETE FROM nerve_rings WHERE uid = $1 AND timepoint = $2"

	_, err := r.DB.Exec(ctx, query, uid, timepoint)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresNerveRingRepository) IngestNerveRing(ctx context.Context, nerveRing domain.NerveRing, skipExisting bool, force bool) (bool, error) {
	exists, err := r.NerveRingExists(ctx, nerveRing.Timepoint)

	if err != nil {
		return false, err
	}

	if skipExisting && exists {
		return true, nil
	}

	if force && exists {
		err = r.DeleteNerveRing(ctx, nerveRing.UID, nerveRing.Timepoint)
		if err != nil {
			return false, err
		}
	}

	err = r.CreateNerveRing(ctx, nerveRing)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *PostgresNerveRingRepository) TruncateNerveRings(ctx context.Context) error {
	query := "TRUNCATE TABLE nerve_rings RESTART IDENTITY CASCADE"

	_, err := r.DB.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
