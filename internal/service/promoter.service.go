package service

import (
	"context"

	"neuroscan/internal/domain"
	"neuroscan/internal/repository"
)

type PromoterService interface {
	GetPromoterByUID(ctx context.Context, uid string) (domain.Promoter, error)
	PromoterExists(ctx context.Context, uid string) (bool, error)
	SearchPromoters(ctx context.Context, query domain.APIV1Request) ([]domain.Promoter, error)
	CountPromoters(ctx context.Context, query domain.APIV1Request) (int, error)
	CreatePromoter(ctx context.Context, promoter domain.Promoter) error
	IngestPromoter(ctx context.Context, promoter domain.Promoter, skipExisting bool, force bool) (bool, error)
	TruncatePromoters(ctx context.Context) error
}

type promoterService struct {
	repo repository.PromoterRepository
}

func NewPromoterService(repo repository.PromoterRepository) PromoterService {
	return &promoterService{
		repo: repo,
	}
}

func (s *promoterService) GetPromoterByUID(ctx context.Context, uid string) (domain.Promoter, error) {
	return s.repo.GetPromoterByUID(ctx, uid)
}

func (s *promoterService) PromoterExists(ctx context.Context, uid string) (bool, error) {
	return s.repo.PromoterExists(ctx, uid)
}

func (s *promoterService) SearchPromoters(ctx context.Context, query domain.APIV1Request) ([]domain.Promoter, error) {
	return s.repo.SearchPromoters(ctx, query)
}

func (s *promoterService) CountPromoters(ctx context.Context, query domain.APIV1Request) (int, error) {
	return s.repo.CountPromoters(ctx, query)
}

func (s *promoterService) CreatePromoter(ctx context.Context, promoter domain.Promoter) error {
	return s.repo.CreatePromoter(ctx, promoter)
}

func (s *promoterService) IngestPromoter(ctx context.Context, promoter domain.Promoter, skipExisting bool, force bool) (bool, error) {
	return s.repo.IngestPromoter(ctx, promoter, skipExisting, force)
}

func (s *promoterService) TruncatePromoters(ctx context.Context) error {
	return s.repo.TruncatePromoters(ctx)
}
