package domain

type DevelopmentalStage struct {
	ID         int    `json:"id"`
	UID        string `json:"uid"`
	Begin      int    `json:"begin"`
	End        int    `json:"end"`
	Order      int    `json:"order"`
	PromoterDB bool   `json:"promoterDB"`
	Timepoints []int  `json:"timepoints"`
}
