package domain

import (
	"errors"

	"neuroscan/internal/toolshed"
)

const NeuronULIDPrefix = "neu"

type Neuron struct {
	ID        int            `json:"-"`
	UID       string         `json:"uid"`
	ULID      string         `json:"id"`
	Timepoint int            `json:"timepoint"`
	Filename  string         `json:"filename"`
	Color     toolshed.Color `json:"color"`
}

func (n *Neuron) Parse(filePath string) error {
	fileMetas, err := toolshed.FilePathParse(filePath)

	if err != nil {
		return errors.New("error parsing neuron file path: " + err.Error())
	}

	fileMeta := fileMetas[0]

	n.UID = fileMeta.UID
	n.ULID = toolshed.BuildUID(NeuronULIDPrefix)
	n.Filename = fileMeta.Filename
	n.Timepoint = fileMeta.Timepoint
	n.Color = fileMeta.Color

	return nil
}

func (n *Neuron) Validate() error {
	if n.ID == 0 {
		return errors.New("id is invalid")
	}

	if n.UID == "" {
		return errors.New("uid is required")
	}

	if n.Filename == "" {
		return errors.New("filename is required")
	}

	return nil
}
