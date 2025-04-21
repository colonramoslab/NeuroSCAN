package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"neuroscan/internal/domain"
	"neuroscan/internal/repository"
)

type NeuronService interface {
	GetNeuronByULID(ctx context.Context, id string) (domain.Neuron, error)
	GetNeuronByUID(ctx context.Context, uid string, timepoint int) (domain.Neuron, error)
	NeuronExists(ctx context.Context, uid string, timepoint int) (bool, error)
	SearchNeurons(ctx context.Context, query domain.APIV1Request) ([]domain.Neuron, error)
	CountNeurons(ctx context.Context, query domain.APIV1Request) (int, error)
	CreateNeuron(ctx context.Context, neuron domain.Neuron) error
	UpdateNeuron(ctx context.Context, neuron domain.Neuron) error
	IngestNeuron(ctx context.Context, neuron domain.Neuron, skipExisting bool, force bool) (bool, error)
	TruncateNeurons(ctx context.Context) error
	ParseMeta(ctx context.Context, row []string, timepoint int, dataType string) error
}

type neuronService struct {
	repo repository.NeuronRepository
}

func NewNeuronService(repo repository.NeuronRepository) NeuronService {
	return &neuronService{
		repo: repo,
	}
}

func (s *neuronService) GetNeuronByULID(ctx context.Context, id string) (domain.Neuron, error) {
	return s.repo.GetNeuronByULID(ctx, id)
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

func (s *neuronService) CreateNeuron(ctx context.Context, neuron domain.Neuron) error {
	return s.repo.CreateNeuron(ctx, neuron)
}

func (s *neuronService) UpdateNeuron(ctx context.Context, neuron domain.Neuron) error {
	return s.repo.UpdateNeuron(ctx, neuron)
}

func (s *neuronService) IngestNeuron(ctx context.Context, neuron domain.Neuron, skipExisting bool, force bool) (bool, error) {
	return s.repo.IngestNeuron(ctx, neuron, skipExisting, force)
}

func (s *neuronService) TruncateNeurons(ctx context.Context) error {
	return s.repo.TruncateNeurons(ctx)
}

func (s *neuronService) ParseMeta(ctx context.Context, row []string, timepoint int, dataType string) error {
	uid := row[0]
	val := row[1]
	neuron, err := s.repo.GetNeuronByUID(ctx, uid, timepoint)
	if err != nil {
		errorString := fmt.Sprintf("unable to find neuron from meta uid %s and timepoint %d: %s", uid, timepoint, err.Error())
		return errors.New(errorString)
	}

	neuron.CellStats = &domain.CellStats{}

	switch dataType {
	case "surface_area":
		value, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		neuron.CellStats.SurfaceArea = &value
	case "volume":
		value, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		neuron.CellStats.Volume = &value

	default:
		return errors.New("unknown data type")
	}

	return s.repo.UpdateNeuron(ctx, neuron)
}
