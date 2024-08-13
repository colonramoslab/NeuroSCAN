package models

type DevStage struct {
	ID         *int    `json:"id"`
	UID        *string `json:"uid"`
	Name       *string `json:"name"`
	Begin      *int    `json:"begin"`
	End        *int    `json:"end"`
	Order      *int    `json:"order"`
	PromoterDB *bool   `json:"promoterDB"`
	Timepoints *string `json:"timepoints"`
}
