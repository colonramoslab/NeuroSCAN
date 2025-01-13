package domain

import (
	"context"

	"neuroscan/internal/logging"
	"neuroscan/internal/toolshed"
)

type Neuron struct {
	ID        int    `json:"id"`
	UID       string `json:"uid"`
	Timepoint int    `json:"timepoint"`
	Filename  string `json:"filename"`
	Color     toolshed.Color `json:"color"`
}

func ParseNeuron(ctx context.Context, filePath string) (Neuron, error) {
	logger := logging.FromContext(ctx)

	fileMetas, err := toolshed.FilePathParse(filePath)

	if err != nil {
		logger.Error().Err(err).Msg("Error parsing file path")
		return Neuron{}, err
	}

	fileMeta := fileMetas[0]

	neuron := Neuron{
		UID:       fileMeta.UID,
		Filename:  fileMeta.Filename,
		Timepoint: fileMeta.Timepoint,
		Color:     fileMeta.Color,
	}

	return neuron, nil
}