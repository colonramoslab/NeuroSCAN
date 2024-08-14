package handlers

import (
	"context"

	"neuroscan/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GetDevStageByUID gets the developmental stage by UID and returns it
func GetDevStageByUID(ctx context.Context, conn *pgxpool.Pool, uid *string) (models.DevStage, error) {
	var devStage models.DevStage

	err := conn.QueryRow(ctx, "SELECT id, uid, name, \"begin\", \"end\", \"order\", \"promoterDB\", timepoints FROM developmental_stages WHERE uid = $1", uid).Scan(&devStage.ID, &devStage.UID, &devStage.Name, &devStage.Begin, &devStage.End, &devStage.Order, &devStage.PromoterDB, &devStage.Timepoints)
	if err != nil {
		return devStage, err
	}

	return devStage, nil
}

// GetDevStageByTimepoint gets the developmental stage by timepoint and returns it
func GetDevStageByTimepoint(ctx context.Context, conn *pgxpool.Pool, timepoint int) (models.DevStage, error) {
	var devStage models.DevStage

	err := conn.QueryRow(ctx, "SELECT id, uid, name, \"begin\", \"end\", \"order\", \"promoterDB\", timepoints FROM developmental_stages WHERE \"begin\" <= $1 AND \"end\" >= $1", timepoint).Scan(&devStage.ID, &devStage.UID, &devStage.Name, &devStage.Begin, &devStage.End, &devStage.Order, &devStage.PromoterDB, &devStage.Timepoints)
	if err != nil {
		return devStage, err
	}

	return devStage, nil
}
