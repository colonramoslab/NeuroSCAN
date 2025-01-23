package domain

import (
	"errors"
	"fmt"
	"strconv"

	"neuroscan/internal/toolshed"
)

const PromoterULIDPrefix = "prmtr"

type Promoter struct {
	ID                        int    `json:"-"`
	ULID                      string `json:"id"`
	UID                       string `json:"uid"`
	Wormbase                  string `json:"wormbase"`
	CellularExpressionPattern string `json:"cellular_expression_pattern"`
	TimepointStart            int    `json:"timepoint_start"`
	TimepointEnd              int    `json:"timepoint_end"`
	CellsByLineaging          string `json:"cells_by_lineaging"`
	ExpressionPatterns        string `json:"expression_patterns"`
	Information               string `json:"information"`
	OtherCells                string `json:"other_cells"`
}

func (p *Promoter) ParseCSV(row []string) error {

	if len(row) != 8 {
		return errors.New("promoter file is invalid")
	}

	timepointStart, _ := strconv.Atoi(row[3])
	timepointEnd, _ := strconv.Atoi(row[4])

	p.UID = row[0]
	p.ULID = toolshed.CreateULID(PromoterULIDPrefix)
	p.Wormbase = row[1]
	p.TimepointStart = timepointStart
	p.TimepointEnd = timepointEnd
	p.CellsByLineaging = row[4]
	p.ExpressionPatterns = row[5]
	p.Information = row[6]
	p.OtherCells = row[7]

	err := p.Validate()
	if err != nil {
		return fmt.Errorf("error validating promoter: %w", err)
	}

	return nil
}

func (p *Promoter) Validate() error {
	if p.UID == "" {
		return errors.New("uid is required")
	}

	return nil
}
