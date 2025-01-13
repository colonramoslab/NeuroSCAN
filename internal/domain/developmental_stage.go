package domain

type DevStage struct {
	id         int
	uid        string
	name       string
	begin      int
	end        int
	order      int
	promoterDB bool
	timepoints []int
}