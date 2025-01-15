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
