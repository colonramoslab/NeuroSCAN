package ingestion

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
)

type DevStage struct {
	id         sql.NullInt64
	uid        string
	name       string
	begin      int
	end        int
	order      int
	promoterDB sql.NullBool
	timepoints string
}

// GetDevStagesAll gets all developmental stages and returns them
func (n *Neuroscan) GetDevStagesAll() ([]DevStage, error) {
	var devStages []DevStage

	rows, err := n.connPool.Query(n.context, "SELECT id, uid, name, \"begin\", \"end\", \"order\", \"promoterDB\", timepoints FROM developmental_stages")
	if err != nil {
		return devStages, err
	}

	defer rows.Close()

	for rows.Next() {
		var devStage DevStage

		err := rows.Scan(&devStage.id, &devStage.uid, &devStage.name, &devStage.begin, &devStage.end, &devStage.order, &devStage.promoterDB, &devStage.timepoints)
		if err != nil {
			return devStages, err
		}

		devStages = append(devStages, devStage)
	}

	return devStages, nil
}

// GetDevStageByUID gets the developmental stage by UID and returns it
func (n *Neuroscan) GetDevStageByUID(uid string) (DevStage, error) {
	var devStage DevStage

	err := n.connPool.QueryRow(n.context, "SELECT id, uid, name, \"begin\", \"end\", \"order\", \"promoterDB\", timepoints FROM developmental_stages WHERE uid = $1", uid).Scan(&devStage.id, &devStage.uid, &devStage.name, &devStage.begin, &devStage.end, &devStage.order, &devStage.promoterDB, &devStage.timepoints)
	if err != nil {
		return devStage, err
	}

	return devStage, nil
}

// GetDevStageByTimepoint gets the developmental stage by timepoint and returns it
func (n *Neuroscan) GetDevStageByTimepoint(timepoint int) DevStage {
	var foundDevStage DevStage

	// search the DevStages on the Neuroscan struct where the comma separated timepoint list contains the timepoint
	for _, devStage := range n.DevStages {
		if devStage.timepoints != "" {
			for _, tp := range strings.Split(devStage.timepoints, ",") {
				if tp == strconv.Itoa(timepoint) {
					foundDevStage = devStage
					break
				}
			}
		}
	}

	return foundDevStage
}

//// GetDevStageByTimepoint gets the developmental stage by timepoint and returns it
//func (n *Neuroscan) GetDevStageByTimepoint(timepoint int) (DevStage, error) {
//	var devStage DevStage
//
//	err := n.connPool.QueryRow(n.context, "SELECT id, uid, name, \"begin\", \"end\", \"order\", \"promoterDB\", timepoints FROM developmental_stages WHERE \"begin\" <= $1 AND \"end\" >= $1", timepoint).Scan(&devStage.id, &devStage.uid, &devStage.name, &devStage.begin, &devStage.end, &devStage.order, &devStage.promoterDB, &devStage.timepoints)
//	if err != nil {
//		return devStage, err
//	}
//
//	return devStage, nil
//}
