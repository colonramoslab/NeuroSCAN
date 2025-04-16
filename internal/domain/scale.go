package domain

import (
	"errors"

	"neuroscan/internal/toolshed"
)

const ScaleULIDPrefix = "scl"

type Scale struct {
	ID        int            `json:"-"`
	ULID      string         `json:"id"`
	UID       string         `json:"uid"`
	Timepoint int            `json:"timepoint"`
	Filename  string         `json:"filename"`
	Color     toolshed.Color `json:"color"`
}

func (s *Scale) Parse(filePath string) error {
	fileMetas, err := toolshed.FilePathParse(filePath)
	if err != nil {
		return errors.New("error parsing scale file path: " + err.Error())
	}

	fileMeta := fileMetas[0]
	ulid := toolshed.CreateULID(ScaleULIDPrefix)

	s.UID = fileMeta.UID
	s.ULID = ulid
	s.Filename = fileMeta.Filename
	s.Timepoint = fileMeta.Timepoint
	s.Color = fileMeta.Color

	return nil
}

func (s *Scale) Validate() error {
	if s.ID == 0 {
		return errors.New("id is invalid")
	}

	if s.UID == "" {
		return errors.New("uid is required")
	}

	if s.Filename == "" {
		return errors.New("filename is required")
	}

	return nil
}
