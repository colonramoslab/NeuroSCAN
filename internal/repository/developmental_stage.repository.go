package repository

import (
	"context"

	"neuroscan/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DevelopmentalStageRepository interface {
	DevelopmentalStageExists(ctx context.Context, uid string) (bool, error)
	SearchDevelopmentalStages(ctx context.Context, query domain.APIV1Request) ([]domain.DevelopmentalStage, error)
	CountDevelopmentalStages(ctx context.Context, query domain.APIV1Request) (int, error)
	CreateDevelopmentalStage(ctx context.Context, developmentalStage domain.DevelopmentalStage) error
}

type PostgresDevelopmentalStageRepository struct {
	DB *pgxpool.Pool
}

func NewPostgresDevelopmentalStageRepository(db *pgxpool.Pool) *PostgresDevelopmentalStageRepository {
	return &PostgresDevelopmentalStageRepository{DB: db}
}

func (r *PostgresDevelopmentalStageRepository) DevelopmentalStageExists(ctx context.Context, uid string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM developmental_stages WHERE uid = $1)"

	var exists bool
	err := r.DB.QueryRow(ctx, query, uid).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PostgresDevelopmentalStageRepository) SearchDevelopmentalStages(ctx context.Context, query domain.APIV1Request) ([]domain.DevelopmentalStage, error) {
	q := "SELECT id, uid, begin, end, order, promoter_db, timepoints FROM developmental_stages "

	parsedQuery, args := query.ToPostgresQuery()

	q += parsedQuery

	rows, _ := r.DB.Query(ctx, q, args...)

	developmentalStages, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.DevelopmentalStage])
	if err != nil {
		return nil, err
	}

	return developmentalStages, nil
}

func (r *PostgresDevelopmentalStageRepository) CountDevelopmentalStages(ctx context.Context, query domain.APIV1Request) (int, error) {
	var count int

	q := "SELECT COUNT(*) FROM developmental_stages"

	parsedQuery, args := query.ToPostgresQuery()

	q += parsedQuery

	err := r.DB.QueryRow(ctx, q, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresDevelopmentalStageRepository) CreateDevelopmentalStage(ctx context.Context, developmentalStage domain.DevelopmentalStage) error {
	query := "INSERT INTO developmental_stages (uid, begin, end, order, promoter_db, timepoints) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING"

	_, err := r.DB.Exec(ctx, query, developmentalStage.UID, developmentalStage.Begin, developmentalStage.End, developmentalStage.Order, developmentalStage.PromoterDB, developmentalStage.Timepoints)
	if err != nil {
		return err
	}

	return nil
}
