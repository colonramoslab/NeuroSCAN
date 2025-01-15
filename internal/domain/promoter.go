package domain

type Promoter struct {
	ID        int    `json:"id"`
	UID       string `json:"uid"`
	Wormbase  string `json:"wormbase"`
	CellularExpressionPattern string `json:"cellular_expression_pattern"`
	TimepointStart int `json:"timepoint_start"`
	TimepointEnd int `json:"timepoint_end"`
	CellsByLineaging string `json:"cells_by_lineaging"`
	ExpressionPatterns string `json:"expression_patterns"`
	Information string `json:"information"`
	OtherCells string `json:"other_cells"`
}
