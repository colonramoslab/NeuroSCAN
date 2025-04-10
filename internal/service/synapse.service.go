package service

import (
	"context"

	"neuroscan/internal/domain"
	"neuroscan/internal/repository"
)

type SynapseService interface {
	GetSynapseByULID(ctx context.Context, id string) (domain.Synapse, error)
	GetSynapseByUID(ctx context.Context, uid string, timepoint int) (domain.Synapse, error)
	SynapseExists(ctx context.Context, uid string, timepoint int) (bool, error)
	SearchSynapses(ctx context.Context, query domain.APIV1Request) ([]domain.Synapse, error)
	CountSynapses(ctx context.Context, query domain.APIV1Request) (int, error)
	CreateSynapse(ctx context.Context, synapse domain.Synapse) error
	IngestSynapse(ctx context.Context, synapse domain.Synapse, skipExisting bool, force bool) (bool, error)
	TruncateSynapses(ctx context.Context) error
}

type synapseService struct {
	repo repository.SynapseRepository
}

func NewSynapseService(repo repository.SynapseRepository) SynapseService {
	return &synapseService{
		repo: repo,
	}
}

func (s *synapseService) GetSynapseByULID(ctx context.Context, id string) (domain.Synapse, error) {
	return s.repo.GetSynapseByULID(ctx, id)
}

func (s *synapseService) GetSynapseByUID(ctx context.Context, uid string, timepoint int) (domain.Synapse, error) {
	return s.repo.GetSynapseByUID(ctx, uid, timepoint)
}

func (s *synapseService) SynapseExists(ctx context.Context, uid string, timepoint int) (bool, error) {
	return s.repo.SynapseExists(ctx, uid, timepoint)
}

func (s *synapseService) SearchSynapses(ctx context.Context, query domain.APIV1Request) ([]domain.Synapse, error) {
	return s.repo.SearchSynapses(ctx, query)
}

func (s *synapseService) CountSynapses(ctx context.Context, query domain.APIV1Request) (int, error) {
	return s.repo.CountSynapses(ctx, query)
}

func (s *synapseService) CreateSynapse(ctx context.Context, synapse domain.Synapse) error {
	return s.repo.CreateSynapse(ctx, synapse)
}

func (s *synapseService) IngestSynapse(ctx context.Context, synapse domain.Synapse, skipExisting bool, force bool) (bool, error) {
	return s.repo.IngestSynapse(ctx, synapse, skipExisting, force)
}

func (s *synapseService) TruncateSynapses(ctx context.Context) error {
	return s.repo.TruncateSynapses(ctx)
}
