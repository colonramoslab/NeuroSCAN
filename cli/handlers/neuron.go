package handlers

import (
	"context"
	"errors"

	"neuroscan/models"
	"neuroscan/toolshed"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetNeuron gets the neuron by UID and returns it, taking an optional timepoint and developmental stage
// the timepoint defaults to 0, and the developmental stage defaults to 0
func GetNeuron(ctx context.Context, conn *pgxpool.Pool, uid *string, timepoint *int, devStage *int) (models.Neuron, error) {
	log.Debug("Getting neuron", "uid", uid, "timepoint", timepoint, "devStage", devStage)
	var neuron models.Neuron

	query := "SELECT neurons.id, uid, embryonic, filename, file_hash, developmental_stage, timepoint FROM neurons WHERE uid = $1 AND timepoint = $2 AND developmental_stage = $3"

	args := []interface{}{uid, timepoint, devStage}

	err := conn.QueryRow(ctx, query, args...).Scan(&neuron.ID, &neuron.UID, &neuron.Embryonic, &neuron.Filename, &neuron.FileHash, &neuron.DevelopmentalStage, &neuron.Timepoint)
	if err != nil {
		log.Error("Error getting neuron: ", "uid", uid, "timepoint", timepoint, "err", err)
		return neuron, err
	}

	return neuron, nil
}

// NeuronExists checks if a neuron exists by the given uid and timepoint
func NeuronExists(ctx context.Context, conn *pgxpool.Pool, uid *string, timepoint *int) (bool, error) {
	var exists bool

	err := conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM neurons WHERE uid = $1 AND timepoint = $2)", uid, timepoint).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// DeleteNeuron deletes a neuron record from the database by the given uid and timepoint
func DeleteNeuron(ctx context.Context, conn *pgxpool.Pool, uid *string, timepoint *int) error {
	_, err := conn.Exec(ctx, "DELETE FROM neurons WHERE uid = $1 AND timepoint = $2", uid, timepoint)
	if err != nil {
		return err
	}

	return nil
}

// CreateNeuron creates a new neuron record in the database
func CreateNeuron(ctx context.Context, conn *pgxpool.Pool, uid *string, embryonic *bool, filename *string, fileHash *string, developmentalStage *int, timepoint *int) error {
	exists, err := NeuronExists(ctx, conn, uid, timepoint)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("neuron already exists")
	}

	_, err = conn.Exec(ctx, "INSERT INTO neurons (uid, embryonic, filename, file_hash, timepoint, developmental_stage) VALUES ($1, $2, $3, $4, $5, $6)", uid, embryonic, filename, fileHash, timepoint, developmentalStage)
	if err != nil {
		return err
	}

	newNeuron, err := GetNeuron(ctx, conn, uid, timepoint, developmentalStage)

	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, "INSERT INTO neurons__developmental_stages (neuron_id, \"developmental-stage_id\") VALUES ($1, $2)", newNeuron.ID, developmentalStage)
	if err != nil {
		return err
	}

	return nil
}

// writeNeuronToDB writes the neuron to the database
func writeNeuronToDB(ctx context.Context, conn *pgxpool.Pool, neuron models.Neuron) {
	exists, err := NeuronExists(ctx, conn, neuron.UID, neuron.Timepoint)

	if err != nil {
		log.Error("Error checking if neuron exists", "err", err)
	}

	// if the neurons exists and we skip existing, return
	//if n.skipExisting && exists {
	//	log.Debug("Neuron exists, skipping", "uid", neuron.UID)
	//	return
	//}

	if exists {
		err = DeleteNeuron(ctx, conn, neuron.UID, neuron.Timepoint)
		if err != nil {
			log.Error("Error deleting existing neuron", "err", err)
			return
		}
	}

	err = CreateNeuron(ctx, conn, neuron.UID, neuron.Embryonic, neuron.Filename, neuron.FileHash, neuron.DevelopmentalStage.ID, neuron.Timepoint)
	if err != nil {
		log.Error("Error inserting new neuron record", "err", err)
		return
	}

	//n.IncrementNeuron()

	log.Debug("Neuron written to database", "uid", neuron.UID, "timepoint", neuron.Timepoint)
}

// parseNeuron takes filename and returns a neuron object
func parseNeuron(ctx context.Context, conn *pgxpool.Pool, filePath string) (models.Neuron, error) {
	fileMetas, err := toolshed.FilePathParse(filePath)

	if err != nil {
		log.Error("Error parsing file path", "err", err)
		return models.Neuron{}, err
	}

	fileMeta := fileMetas[0]

	devStage, err := GetDevStageByUID(ctx, conn, fileMeta.DevelopmentalStage)

	if err != nil {
		log.Error("Error getting developmental stage", "err", err)
		return models.Neuron{}, err
	}

	var embryonic bool
	embryonic = false

	neuron := models.Neuron{
		UID:                fileMeta.UID,
		Embryonic:          &embryonic,
		Filename:           fileMeta.Filename,
		FileHash:           fileMeta.Filehash,
		DevelopmentalStage: &devStage,
		Timepoint:          fileMeta.Timepoint,
	}

	return neuron, nil
}

// ProcessNeuron processes a neuron file from the path, writing it to the database
func ProcessNeuron(ctx context.Context, conn *pgxpool.Pool, filePath string) {
	neuron, err := parseNeuron(ctx, conn, filePath)
	if err != nil {
		log.Error("Error parsing neuron", "err", err)
		return
	}

	writeNeuronToDB(ctx, conn, neuron)
}
