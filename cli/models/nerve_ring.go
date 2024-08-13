package models

type NerveRing struct {
	ID                 *int      `json:"id"`
	UID                *string   `json:"uid"`
	DevelopmentalStage *DevStage `json:"developmental_stage"`
	Timepoint          *int      `json:"timepoint"`
	Filename           *string   `json:"filename"`
	FileHash           *string   `json:"file_hash"`
}
