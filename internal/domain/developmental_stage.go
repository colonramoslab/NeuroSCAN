package domain

const DevelopmentalStageULIDPrefix = "devstg"

type DevelopmentalStage struct {
	ID         int    `json:"-"`
	UID        string `json:"uid"`
	ULID       string `json:"id"`
	Begin      int    `json:"begin"`
	End        int    `json:"end"`
	Order      int    `json:"order"`
	PromoterDB *bool  `json:"promoterDB"`
	Timepoints []int  `json:"timepoints"`
}
