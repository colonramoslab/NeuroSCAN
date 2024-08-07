package ingestion

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
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
func (n *Neuroscan) GetDevStageByTimepoint(timepoint int) (DevStage, error) {
	var devStage DevStage

	err := n.connPool.QueryRow(n.context, "SELECT id, uid, name, \"begin\", \"end\", \"order\", \"promoterDB\", timepoints FROM developmental_stages WHERE \"begin\" <= $1 AND \"end\" >= $1", timepoint).Scan(&devStage.id, &devStage.uid, &devStage.name, &devStage.begin, &devStage.end, &devStage.order, &devStage.promoterDB, &devStage.timepoints)
	if err != nil {
		return devStage, err
	}

	return devStage, nil
}
