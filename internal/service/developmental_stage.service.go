package service

import (
	"context"

	"neuroscan/internal/domain"
	"neuroscan/internal/repository"
)

type DevelopmentalStageService interface {
	SearchDevelopmentalStages(ctx context.Context, query domain.APIV1Request) ([]domain.DevelopmentalStage, error)
	CountDevelopmentalStages(ctx context.Context, query domain.APIV1Request) (int, error)
	CreateDevelopmentalStage(ctx context.Context, developmentalStage domain.DevelopmentalStage) error
	TruncateDevelopmentalStages(ctx context.Context) error
}

type developmentalStageService struct {
	repo repository.DevelopmentalStageRepository
}

func NewDevelopmentalStageService(repo repository.DevelopmentalStageRepository) DevelopmentalStageService {
	return &developmentalStageService{repo: repo}
}

func (s *developmentalStageService) SearchDevelopmentalStages(ctx context.Context, query domain.APIV1Request) ([]domain.DevelopmentalStage, error) {
	return s.repo.SearchDevelopmentalStages(ctx, query)
}

func (s *developmentalStageService) CountDevelopmentalStages(ctx context.Context, query domain.APIV1Request) (int, error) {
	return s.repo.CountDevelopmentalStages(ctx, query)
}

func (s *developmentalStageService) CreateDevelopmentalStage(ctx context.Context, developmentalStage domain.DevelopmentalStage) error {
	return s.repo.CreateDevelopmentalStage(ctx, developmentalStage)
}

func (s *developmentalStageService) TruncateDevelopmentalStages(ctx context.Context) error {
	return s.repo.TruncateDevelopmentalStages(ctx)
}
