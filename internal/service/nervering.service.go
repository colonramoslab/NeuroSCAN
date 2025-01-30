package service

import (
	"context"

	"neuroscan/internal/domain"
	"neuroscan/internal/repository"
)

type NerveRingService interface {
	GetNerveRingByTimepoint(ctx context.Context, timepoint int) (domain.NerveRing, error)
	NerveRingExists(ctx context.Context, timepoint int) (bool, error)
	CreateNerveRing(ctx context.Context, nervering domain.NerveRing) error
	IngestNerveRing(ctx context.Context, nervering domain.NerveRing, skipExisting bool, force bool) (bool, error)
	TruncateNerveRings(ctx context.Context) error
}

type nerveringService struct {
	repo repository.NerveRingRepository
}

func NewNerveRingService(repo repository.NerveRingRepository) NerveRingService {
	return &nerveringService{
		repo: repo,
	}
}

func (s *nerveringService) GetNerveRingByTimepoint(ctx context.Context, timepoint int) (domain.NerveRing, error) {
	return s.repo.GetNerveRingByTimepoint(ctx, timepoint)
}

func (s *nerveringService) NerveRingExists(ctx context.Context, timepoint int) (bool, error) {
	return s.repo.NerveRingExists(ctx, timepoint)
}

func (s *nerveringService) CreateNerveRing(ctx context.Context, nervering domain.NerveRing) error {
	return s.repo.CreateNerveRing(ctx, nervering)
}

func (s *nerveringService) IngestNerveRing(ctx context.Context, nervering domain.NerveRing, skipExisting bool, force bool) (bool, error) {
	return s.repo.IngestNerveRing(ctx, nervering, skipExisting, force)
}

func (s *nerveringService) TruncateNerveRings(ctx context.Context) error {
	return s.repo.TruncateNerveRings(ctx)
}
