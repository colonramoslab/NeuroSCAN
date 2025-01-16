package domain

import (
	"neuroscan/internal/toolshed"
)

type CphateStructure []struct {
	I       int            `json:"i"`
	C       int            `json:"c"`
	Neurons []string       `json:"neurons"`
	ObjFile string         `json:"objFile"`
	Color   toolshed.Color `json:"color"`
}

type Cphate struct {
	ID        int             `json:"id"`
	UID       string          `json:"uid"`
	Timepoint int             `json:"timepoint"`
	Structure CphateStructure `json:"structure"`
	Filename  string          `json:"filename"`
}
