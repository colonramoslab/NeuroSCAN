package domain

const PromoterULIDPrefix = "prmtr"

type Promoter struct {
	ID                        int    `json:"-"`
	UID                       string `json:"uid"`
	ULID                      string `json:"id"`
	Wormbase                  string `json:"wormbase"`
	CellularExpressionPattern string `json:"cellular_expression_pattern"`
	TimepointStart            int    `json:"timepoint_start"`
	TimepointEnd              int    `json:"timepoint_end"`
	CellsByLineaging          string `json:"cells_by_lineaging"`
	ExpressionPatterns        string `json:"expression_patterns"`
	Information               string `json:"information"`
	OtherCells                string `json:"other_cells"`
}
