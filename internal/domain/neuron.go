package domain

import (
	"neuroscan/internal/toolshed"
)

type Neuron struct {
	ID        int    `json:"id"`
	UID       string `json:"uid"`
	Timepoint int    `json:"timepoint"`
	Filename  string `json:"filename"`
	Color     toolshed.Color `json:"color"`
}
