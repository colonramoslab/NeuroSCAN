package domain

import (
	"errors"
	"strings"

	"neuroscan/internal/toolshed"
)

type SynapseType string

const (
	SynapseTypeChemical   SynapseType = "chemical"
	SynapseTypeElectrical SynapseType = "electrical"
	SynapseTypeUndefined  SynapseType = "undefined"
	SynapseULIDPrefix                 = "syn"
)

var ValidSynapseType = map[SynapseType]bool{
	SynapseTypeChemical:   true,
	SynapseTypeElectrical: true,
}

type Synapse struct {
	ID           int            `json:"-"`
	ULID         string         `json:"id"`
	UID          string         `json:"uid"`
	Timepoint    int            `json:"timepoint"`
	SynapseType  SynapseType    `json:"type"`
	Filename     string         `json:"filename"`
	Color        toolshed.Color `json:"color"`
	CellStats    *CellStats     `json:"cell_stats"`
	SynapseStats *SynapseStats  `json:"synapse_stats"`
}

type SynapseStats struct {
	TotalTypeCount        *int           `json:"total_type_count"`
	TotalCellSynapseCount *int           `json:"total_cell_synapse_count"`
	Connections           *[]SynapseItem `json:"connections"`
}

type SynapseItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func getSynapseType(uid string) *SynapseType {
	var synapseType SynapseType

	if strings.Contains(uid, "chemical") {
		synapseType = SynapseTypeChemical
	}

	if strings.Contains(uid, "electrical") {
		synapseType = SynapseTypeElectrical
	}

	return &synapseType
}

func (s *Synapse) Parse(filePath string) error {
	fileMetas, err := toolshed.FilePathParse(filePath)
	if err != nil {
		return errors.New("error parsing synapse file path: " + err.Error())
	}

	fileMeta := fileMetas[0]
	ulid := toolshed.CreateULID(SynapseULIDPrefix)
	synapseType := getSynapseType(fileMeta.UID)

	s.UID = fileMeta.UID
	s.ULID = ulid
	s.Filename = fileMeta.Filename
	s.Timepoint = fileMeta.Timepoint
	s.Color = fileMeta.Color
	s.SynapseType = *synapseType

	return nil
}

func (s *Synapse) Validate() error {
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
