package domain

type Ranking struct {
	CellRank         int     `json:"cell_rank"`
	CellTotal        int     `json:"cell_total"`
	CellSAAggregate  float64 `json:"cell_sa_aggregate"`
	BrainRank        int     `json:"brain_rank"`
	BrainTotal       int     `json:"brain_total"`
	BrainSAAggregate float64 `json:"brain_sa_aggregate"`
}
