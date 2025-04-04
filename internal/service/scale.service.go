package service

import (
	"context"

	"neuroscan/internal/domain"
	"neuroscan/internal/repository"
)

type ScaleService interface {
	GetScaleByTimepoint(ctx context.Context, timepoint int) ([]domain.Scale, error)
	ScaleExists(ctx context.Context, timepoint int) (bool, error)
	CreateScale(ctx context.Context, scale domain.Scale) error
	IngestScale(ctx context.Context, scale domain.Scale, skipExisting bool, force bool) (bool, error)
	TruncateScales(ctx context.Context) error
}

type scaleService struct {
	repo repository.ScaleRepository
}

func NewScaleService(repo repository.ScaleRepository) ScaleService {
	return &scaleService{
		repo: repo,
	}
}

func (s *scaleService) GetScaleByTimepoint(ctx context.Context, timepoint int) ([]domain.Scale, error) {
	return s.repo.GetScaleByTimepoint(ctx, timepoint)
}

func (s *scaleService) ScaleExists(ctx context.Context, timepoint int) (bool, error) {
	return s.repo.ScaleExists(ctx, timepoint)
}

func (s *scaleService) CreateScale(ctx context.Context, scale domain.Scale) error {
	return s.repo.CreateScale(ctx, scale)
}

func (s *scaleService) IngestScale(ctx context.Context, scale domain.Scale, skipExisting bool, force bool) (bool, error) {
	return s.repo.IngestScale(ctx, scale, skipExisting, force)
}

func (s *scaleService) TruncateScales(ctx context.Context) error {
	return s.repo.TruncateScales(ctx)
}
