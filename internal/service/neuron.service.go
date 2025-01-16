package service

import (
	"context"
	"errors"

	"neuroscan/internal/domain"
	"neuroscan/internal/repository"
	"neuroscan/internal/toolshed"
)

type NeuronService interface {
	GetNeuronByID(ctx context.Context, uid string, timepoint int) (domain.Neuron, error)
	GetNeuronByUID(ctx context.Context, uid string, timepoint int) (domain.Neuron, error)
	NeuronExists(ctx context.Context, uid string, timepoint int) (bool, error)
	SearchNeurons(ctx context.Context, query domain.APIV1Request) ([]domain.Neuron, error)
	CountNeurons(ctx context.Context, query domain.APIV1Request) (int, error)
	CreateNeuron(ctx context.Context, uid string, filename string, timepoint int, color toolshed.Color) error
	IngestNeuron(ctx context.Context, neuron domain.Neuron, skipExisting bool, force bool) (bool, error)
	ParseNeuron(ctx context.Context, filePath string) (domain.Neuron, error)
}

type neuronService struct {
	repo repository.NeuronRepository
}

func NewNeuronService(repo repository.NeuronRepository) NeuronService {
	return &neuronService{
		repo: repo,
	}
}

func (s *neuronService) GetNeuronByID(ctx context.Context, uid string, timepoint int) (domain.Neuron, error) {
	return s.repo.GetNeuronByID(ctx, uid, timepoint)
}

func (s *neuronService) GetNeuronByUID(ctx context.Context, uid string, timepoint int) (domain.Neuron, error) {
	return s.repo.GetNeuronByUID(ctx, uid, timepoint)
}

func (s *neuronService) NeuronExists(ctx context.Context, uid string, timepoint int) (bool, error) {
	return s.repo.NeuronExists(ctx, uid, timepoint)
}

func (s *neuronService) SearchNeurons(ctx context.Context, query domain.APIV1Request) ([]domain.Neuron, error) {
	return s.repo.SearchNeurons(ctx, query)
}

func (s *neuronService) CountNeurons(ctx context.Context, query domain.APIV1Request) (int, error) {
	return s.repo.CountNeurons(ctx, query)
}

func (s *neuronService) CreateNeuron(ctx context.Context, uid string, filename string, timepoint int, color toolshed.Color) error {
	return s.repo.CreateNeuron(ctx, uid, filename, timepoint, color)
}

func (s *neuronService) IngestNeuron(ctx context.Context, neuron domain.Neuron, skipExisting bool, force bool) (bool, error) {
	return s.repo.IngestNeuron(ctx, neuron, skipExisting, force)
}

func (s *neuronService) ParseNeuron(ctx context.Context, filePath string) (domain.Neuron, error) {
	fileMetas, err := toolshed.FilePathParse(filePath)

	if err != nil {
		return domain.Neuron{}, errors.New("error parsing neuron file path: " + err.Error())
	}

	fileMeta := fileMetas[0]

	neuron := domain.Neuron{
		UID:       fileMeta.UID,
		Filename:  fileMeta.Filename,
		Timepoint: fileMeta.Timepoint,
		Color:     fileMeta.Color,
	}

	return neuron, nil
}
