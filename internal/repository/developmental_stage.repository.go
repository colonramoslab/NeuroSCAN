package repository

import (
	"context"
	"fmt"
	"strings"

	"neuroscan/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DevelopmentalStageRepository interface {
	DevelopmentalStageExists(ctx context.Context, uid string) (bool, error)
	SearchDevelopmentalStages(ctx context.Context, query domain.APIV1Request) ([]domain.DevelopmentalStage, error)
	CountDevelopmentalStages(ctx context.Context, query domain.APIV1Request) (int, error)
	IngestDevelopmentalStage(ctx context.Context, devStage domain.DevelopmentalStage, skipExisting bool, force bool) (bool, error)
	DeleteDevelopmentalStage(ctx context.Context, developmentalStage domain.DevelopmentalStage) error
	CreateDevelopmentalStage(ctx context.Context, developmentalStage domain.DevelopmentalStage) error
	TruncateDevelopmentalStages(ctx context.Context) error
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
	q := `SELECT id, uid, ulid, begin, "end", "order", promoter_db, timepoints FROM developmental_stages `

	parsedQuery, args := r.ParseDevelopmentalStageAPIV1Request(ctx, query)

	q += parsedQuery

	rows, _ := r.DB.Query(ctx, q, args...)

	var developmentalStages []domain.DevelopmentalStage

	for rows.Next() {
		var developmentalStage domain.DevelopmentalStage

		err := rows.Scan(&developmentalStage.ID, &developmentalStage.UID, &developmentalStage.ULID, &developmentalStage.Begin, &developmentalStage.End, &developmentalStage.Order, &developmentalStage.PromoterDB, &developmentalStage.Timepoints)
		if err != nil {
			return nil, err
		}

		developmentalStages = append(developmentalStages, developmentalStage)
	}

	return developmentalStages, nil
}

func (r *PostgresDevelopmentalStageRepository) CountDevelopmentalStages(ctx context.Context, query domain.APIV1Request) (int, error) {
	var count int

	q := "SELECT COUNT(*) FROM developmental_stages"

	parsedQuery, args := r.ParseDevelopmentalStageAPIV1Request(ctx, query)

	q += parsedQuery

	err := r.DB.QueryRow(ctx, q, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresDevelopmentalStageRepository) CreateDevelopmentalStage(ctx context.Context, developmentalStage domain.DevelopmentalStage) error {
	query := `INSERT INTO developmental_stages (uid, ulid, begin, "end", "order", promoter_db, timepoints) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING`

	_, err := r.DB.Exec(ctx, query, developmentalStage.UID, developmentalStage.ULID, developmentalStage.Begin, developmentalStage.End, developmentalStage.Order, developmentalStage.PromoterDB, developmentalStage.Timepoints)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresDevelopmentalStageRepository) DeleteDevelopmentalStage(ctx context.Context, developmentalStage domain.DevelopmentalStage) error {
	query := `DELETE FROM developmental_stages WHERE uid = $1`

	_, err := r.DB.Exec(ctx, query, developmentalStage.UID)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresDevelopmentalStageRepository) IngestDevelopmentalStage(ctx context.Context, devStage domain.DevelopmentalStage, skipExisting bool, force bool) (bool, error) {
	exists, err := r.DevelopmentalStageExists(ctx, devStage.UID)
	if err != nil {
		return false, err
	}

	if skipExisting && exists {
		return true, nil
	}

	if force && exists {
		err := r.DeleteDevelopmentalStage(ctx, devStage)
		if err != nil {
			return false, err
		}
	}

	err = r.CreateDevelopmentalStage(ctx, devStage)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *PostgresDevelopmentalStageRepository) TruncateDevelopmentalStages(ctx context.Context) error {
	query := "TRUNCATE TABLE developmental_stages RESTART IDENTITY CASCADE"

	_, err := r.DB.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresDevelopmentalStageRepository) ParseDevelopmentalStageAPIV1Request(ctx context.Context, req domain.APIV1Request) (string, []any) {
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

	query := strings.Join(queryParts, " AND ")

	return query, args
}
