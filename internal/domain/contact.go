package domain

import (
	"errors"

	"neuroscan/internal/toolshed"
)

const ContactULIDPrefix = "cntct"

type Contact struct {
	ID        int            `json:"-"`
	ULID      string         `json:"id"`
	UID       string         `json:"uid"`
	Timepoint int            `json:"timepoint"`
	Filename  string         `json:"filename"`
	Color     toolshed.Color `json:"color"`
}

func (c *Contact) Parse(filePath string) error {
	fileMetas, err := toolshed.FilePathParse(filePath)

	if err != nil {
		return errors.New("error parsing contact file path: " + err.Error())
	}

	ulid := toolshed.BuildUID(ContactULIDPrefix)
	fileMeta := fileMetas[0]

	c.UID = fileMeta.UID
	c.ULID = ulid
	c.Filename = fileMeta.Filename
	c.Timepoint = fileMeta.Timepoint
	c.Color = fileMeta.Color

	return nil
}

func (c *Contact) Validate() error {
	if c.ID == 0 {
		return errors.New("id is invalid")
	}

	if c.UID == "" {
		return errors.New("uid is required")
	}

	if c.Filename == "" {
		return errors.New("filename is required")
	}

	return nil
}
