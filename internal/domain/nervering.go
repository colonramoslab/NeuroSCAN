package domain

import (
	"errors"

	"neuroscan/internal/toolshed"
)

type NerveRing struct {
	ID        int            `json:"id"`
	UID       string         `json:"uid"`
	Timepoint int            `json:"timepoint"`
	Filename  string         `json:"filename"`
	Color     toolshed.Color `json:"color"`
}

func (n *NerveRing) Parse(filePath string) error {
	fileMetas, err := toolshed.FilePathParse(filePath)

	if err != nil {
		return errors.New("error parsing nerve ring file path: " + err.Error())
	}

	fileMeta := fileMetas[0]

	n.UID = fileMeta.UID
	n.Filename = fileMeta.Filename
	n.Timepoint = fileMeta.Timepoint
	n.Color = fileMeta.Color

	return nil
}

func (n *NerveRing) Validate() error {
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

