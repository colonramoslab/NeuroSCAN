package domain

import (
	"neuroscan/internal/toolshed"
)

type SynapseType string

const (
	SynapseTypeChemical SynapseType = "chemical"
	SynapseTypeElectrical SynapseType = "electrical"
)

type Synapse struct {
	ID        int    `json:"id"`
	UID       string `json:"uid"`
	Timepoint int    `json:"timepoint"`
	SynapseType string `json:"type"`
	Filename   string `json:"filename"`
	Color      toolshed.Color  `json:"color"`
}
