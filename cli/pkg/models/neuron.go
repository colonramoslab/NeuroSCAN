package models

type Neuron struct {
	ID                 *int      `json:"id"`
	UID                *string   `json:"uid"`
	Embryonic          *bool     `json:"embryonic"`
	Filename           *string   `json:"filename"`
	FileHash           *string   `json:"file_hash"`
	DevelopmentalStage *DevStage `json:"developmental_stage"`
	Timepoint          *int      `json:"timepoint"`
}
