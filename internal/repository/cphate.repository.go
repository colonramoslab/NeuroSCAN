package repository

import (
	"context"
	"errors"
	"fmt"

	"neuroscan/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CphateRepository interface {
	GetCphateByTimepoint(ctx context.Context, timepoint int) (domain.Cphate, error)
	CphateExists(ctx context.Context, timepoint int) (bool, error)
	CreateCphate(ctx context.Context, cphate domain.Cphate) error
	DeleteCphate(ctx context.Context, timepoint int) error
	IngestCphate(ctx context.Context, cphate domain.Cphate, skipExisting bool, force bool) (bool, error)
	TruncateCphates(ctx context.Context) error
}

type PostgresCphateRepository struct {
	DB *pgxpool.Pool
}

func NewPostgresCphateRepository(db *pgxpool.Pool) *PostgresCphateRepository {
	return &PostgresCphateRepository{
		DB: db,
	}
}

func (r *PostgresCphateRepository) GetCphateByTimepoint(ctx context.Context, timepoint int) (domain.Cphate, error) {
	query := "SELECT id, uid, timepoint, structure FROM cphates WHERE timepoint = $1"

	var cphate domain.Cphate
	err := r.DB.QueryRow(ctx, query, timepoint).Scan(&cphate.ID, &cphate.UID, &cphate.Timepoint, &cphate.Structure)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Cphate{}, nil
		}

		return domain.Cphate{}, err
	}

	return cphate, nil
}

func (r *PostgresCphateRepository) CphateExists(ctx context.Context, timepoint int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM cphates WHERE timepoint = $1)"

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

func (r *PostgresCphateRepository) CreateCphate(ctx context.Context, cphate domain.Cphate) error {
	exists, err := r.CphateExists(ctx, cphate.Timepoint)

	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("cphate already exists")
	}

	query := "INSERT INTO cphates (uid, timepoint, structure) VALUES ($1, $2, $3)"

	_, err = r.DB.Exec(ctx, query, cphate.UID, cphate.Timepoint, cphate.Structure)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresCphateRepository) DeleteCphate(ctx context.Context, timepoint int) error {
	query := "DELETE FROM cphates WHERE timepoint = $1"

	_, err := r.DB.Exec(ctx, query, timepoint)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresCphateRepository) IngestCphate(ctx context.Context, cphate domain.Cphate, skipExisting bool, force bool) (bool, error) {
	exists, err := r.CphateExists(ctx, cphate.Timepoint)

	if err != nil {
		return false, err
	}

	if skipExisting && exists {
		return true, nil
	}

	if force && exists {
		err = r.DeleteCphate(ctx, cphate.Timepoint)
		if err != nil {
			return false, err
		}
	}

	err = r.CreateCphate(ctx, cphate)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *PostgresCphateRepository) TruncateCphates(ctx context.Context) error {
	query := "TRUNCATE TABLE cphates RESTART IDENTITY CASCADE"

	_, err := r.DB.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
