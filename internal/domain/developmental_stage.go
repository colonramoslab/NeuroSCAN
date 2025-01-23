package domain

import (
	"errors"
	"fmt"
	"strconv"

	"neuroscan/internal/toolshed"
)

const DevelopmentalStageULIDPrefix = "devstg"

type DevelopmentalStage struct {
	ID         int    `json:"-"`
	ULID       string `json:"id"`
	UID        string `json:"uid"`
	Begin      int    `json:"begin"`
	End        int    `json:"end"`
	Order      int    `json:"order"`
	PromoterDB *bool  `json:"promoterDB"`
	Timepoints []int  `json:"timepoints"`
}

func (ds *DevelopmentalStage) ParseCSV(row []string) error {

	if len(row) != 6 {
		return errors.New("developmental stage file is invalid")
	}

	begin, _ := strconv.Atoi(row[1])
	end, _ := strconv.Atoi(row[2])
	order, _ := strconv.Atoi(row[3])
	promoterDB, _ := strconv.ParseBool(row[4])

	timepoints := toolshed.ParseTimepointIntArray(row[5])

	ds.UID = row[0]
	ds.ULID = toolshed.CreateULID(DevelopmentalStageULIDPrefix)
	ds.Begin = begin
	ds.End = end
	ds.Order = order
	ds.PromoterDB = &promoterDB
	ds.Timepoints = timepoints

	err := ds.Validate()

	if err != nil {
		return fmt.Errorf("developmental stage file is invalid: %w", err)
	}

	return nil
}

func (ds *DevelopmentalStage) Validate() error {
	if ds.UID == "" {
		return errors.New("uid is required")
	}

	return nil
}
