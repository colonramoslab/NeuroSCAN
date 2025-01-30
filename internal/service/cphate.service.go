package service

import (
	"context"

	"neuroscan/internal/domain"
	"neuroscan/internal/repository"
)

type CphateService interface {
	GetCphateByTimepoint(ctx context.Context, timepoint int) (domain.Cphate, error)
	CountCphates(ctx context.Context, timepoint int) (int, error)
	CphateExists(ctx context.Context, timepoint int) (bool, error)
	CreateCphate(ctx context.Context, cphate domain.Cphate) error
	IngestCphate(ctx context.Context, cphate domain.Cphate, skipExisting bool, force bool) (bool, error)
	TruncateCphates(ctx context.Context) error
}

type cphateService struct {
	repo repository.CphateRepository
}

func NewCphateService(repo repository.CphateRepository) CphateService {
	return &cphateService{
		repo: repo,
	}
}

func (s *cphateService) GetCphateByTimepoint(ctx context.Context, timepoint int) (domain.Cphate, error) {
	return s.repo.GetCphateByTimepoint(ctx, timepoint)
}

func (s *cphateService) CountCphates(ctx context.Context, timepoint int) (int, error) {
	return s.repo.CountCphates(ctx, timepoint)
}

func (s *cphateService) CphateExists(ctx context.Context, timepoint int) (bool, error) {
	return s.repo.CphateExists(ctx, timepoint)
}

func (s *cphateService) CreateCphate(ctx context.Context, cphate domain.Cphate) error {
	return s.repo.CreateCphate(ctx, cphate)
}

func (s *cphateService) IngestCphate(ctx context.Context, cphate domain.Cphate, skipExisting bool, force bool) (bool, error) {
	return s.repo.IngestCphate(ctx, cphate, skipExisting, force)
}

func (s *cphateService) TruncateCphates(ctx context.Context) error {
	return s.repo.TruncateCphates(ctx)
}
