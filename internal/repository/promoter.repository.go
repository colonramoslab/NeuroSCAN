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

type PromoterRepository interface {
	GetPromoterByUID(ctx context.Context, uid string) (domain.Promoter, error)
	PromoterExists(ctx context.Context, uid string) (bool, error)
	SearchPromoters(ctx context.Context, query domain.APIV1Request) ([]domain.Promoter, error)
	CountPromoters(ctx context.Context, query domain.APIV1Request) (int, error)
	CreatePromoter(ctx context.Context, promoter domain.Promoter) error
	DeletePromoter(ctx context.Context, uid string) error
	IngestPromoter(ctx context.Context, promoter domain.Promoter, skipExisting bool, force bool) (bool, error)
	TruncatePromoters(ctx context.Context) error
}

type PostgresPromoterRepository struct {
	DB *pgxpool.Pool
}

func NewPostgresPromoterRepository(db *pgxpool.Pool) *PostgresPromoterRepository {
	return &PostgresPromoterRepository{
		DB: db,
	}
}

func (r *PostgresPromoterRepository) GetPromoterByUID(ctx context.Context, uid string) (domain.Promoter, error) {
	query := "SELECT id, uid, ulid, wormbase, cellular_expression_pattern, timepoint_start, timepoint_end, cells_by_lineaging, expression_patterns, information, other_cells FROM promoters WHERE uid = $1"

	var promoter domain.Promoter
	err := r.DB.QueryRow(ctx, query, uid).Scan(&promoter.ID, &promoter.UID, &promoter.ULID, &promoter.Wormbase, &promoter.CellularExpressionPattern, &promoter.TimepointStart, &promoter.TimepointEnd, &promoter.CellsByLineaging, &promoter.ExpressionPatterns, &promoter.Information, &promoter.OtherCells)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Promoter{}, nil
		}

		return domain.Promoter{}, err
	}

	return promoter, nil
}

func (r *PostgresPromoterRepository) PromoterExists(ctx context.Context, uid string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM promoters WHERE uid = $1)"

	var exists bool
	err := r.DB.QueryRow(ctx, query, uid).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return exists, nil
}

func (r *PostgresPromoterRepository) SearchPromoters(ctx context.Context, query domain.APIV1Request) ([]domain.Promoter, error) {
	q := "SELECT id, uid, ulid, wormbase, cellular_expression_pattern, timepoint_start, timepoint_end, cells_by_lineaging, expression_patterns, information, other_cells FROM promoters "

	parsedQuery, args := r.ParsePromoterAPIV1Request(ctx, query)

	q += parsedQuery

	rows, _ := r.DB.Query(ctx, q, args...)

	promoters, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Promoter])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.Promoter{}, nil
		}

		return nil, err
	}

	return promoters, nil
}

func (r *PostgresPromoterRepository) CountPromoters(ctx context.Context, query domain.APIV1Request) (int, error) {
	var count int

	q := "SELECT COUNT(*) FROM promoters "

	parsedQuery, args := r.ParsePromoterAPIV1Request(ctx, query)

	q += parsedQuery

	err := r.DB.QueryRow(ctx, q, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresPromoterRepository) CreatePromoter(ctx context.Context, promoter domain.Promoter) error {
	exists, err := r.PromoterExists(ctx, promoter.UID)

	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("promoter already exists")
	}

	query := "INSERT INTO promoters (uid, ulid, wormbase, cellular_expression_pattern, timepoint_start, timepoint_end, cells_by_lineaging, expression_patterns, information, other_cells) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT DO NOTHING"

	_, err = r.DB.Exec(ctx, query, promoter.UID, promoter.ULID, promoter.Wormbase, promoter.CellularExpressionPattern, promoter.TimepointStart, promoter.TimepointEnd, promoter.CellsByLineaging, promoter.ExpressionPatterns, promoter.Information, promoter.OtherCells)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresPromoterRepository) DeletePromoter(ctx context.Context, uid string) error {
	query := "DELETE FROM promoters WHERE uid = $1"

	_, err := r.DB.Exec(ctx, query, uid)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresPromoterRepository) TruncatePromoters(ctx context.Context) error {
	query := "TRUNCATE TABLE promoters RESTART IDENTITY CASCADE"

	_, err := r.DB.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresPromoterRepository) IngestPromoter(ctx context.Context, promoter domain.Promoter, skipExisting bool, force bool) (bool, error) {
	exists, err := r.PromoterExists(ctx, promoter.UID)

	if err != nil {
		return false, err
	}

	if skipExisting && exists {
		return true, nil
	}

	if force && exists {
		err := r.DeletePromoter(ctx, promoter.UID)
		if err != nil {
			return false, err
		}
	}

	err = r.CreatePromoter(ctx, promoter)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *PostgresPromoterRepository) ParsePromoterAPIV1Request(ctx context.Context, req domain.APIV1Request) (string, []interface{}) {

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
