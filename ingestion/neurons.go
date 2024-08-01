package main

import (
	"database/sql"
	"errors"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
)

type Neuron struct {
	id                 int
	uid                string
	embryonic          bool
	filename           string
	fileHash           string
	developmentalStage sql.NullInt64
	timepoint          int
}

// GetNeuron gets the neuron by UID and returns it, taking an optional timepoint and developmental stage
// the timepoint defaults to 0, and the developmental stage defaults to 0
func (n *Neuroscan) GetNeuron(uid string, timepoint int, devStage sql.NullInt64) (Neuron, error) {
	log.Debug("Getting neuron", "uid", uid, "timepoint", timepoint, "devStage", devStage)
	var neuron Neuron

	db, err := n.ConnectDB()

	if err != nil {
		log.Error("Error connecting to database", "err", err)
		return neuron, err
	}

	defer db.Close()

	query := `
		SELECT neurons.id, uid, embryonic, filename, file_hash, "developmental-stage_id", timepoint
		FROM neurons
		    INNER JOIN neurons__developmental_stages
		        ON neurons.id = neurons__developmental_stages.neuron_id
		WHERE uid = ?
		AND timepoint = ?
	`

	args := []interface{}{uid, timepoint}

	if devStage.Valid {
		query += "AND `developmental-stage_id` = ?"
		args = append(args, devStage)
	}

	err = db.QueryRow(query, args...).Scan(&neuron.id, &neuron.uid, &neuron.embryonic, &neuron.filename, &neuron.fileHash, &neuron.developmentalStage, &neuron.timepoint)
	if err != nil {
		return neuron, err
	}

	return neuron, nil
}

// writeToDB writes the neuron to the database
func (neuron Neuron) writeToDB(n *Neuroscan) {
	exists, err := n.NeuronExists(neuron.uid, neuron.timepoint)

	if err != nil {
		log.Error("Error checking if neuron exists", "err", err)
	}

	db, err := n.ConnectDB()

	if err != nil {
		log.Error("Error connecting to database", "err", err)
		return
	}

	defer db.Close()

	if exists {
		log.Debug("Neuron exists, updating record", "uid", neuron.uid)
		err := n.UpdateNeuron(neuron.uid, neuron.timepoint, neuron.embryonic, neuron.filename, neuron.fileHash, neuron.developmentalStage)

		if err != nil {
			log.Error("Error updating neuron record", "err", err)
			return
		}
	} else {
		log.Debug("Neuron does not exist, inserting new record", "uid", neuron.uid)
		err := n.CreateNeuron(neuron.uid, neuron.embryonic, neuron.filename, neuron.fileHash, neuron.developmentalStage, neuron.timepoint)
		if err != nil {
			log.Error("Error inserting new neuron record", "err", err)
			return
		}
	}

	n.IncrementNeuron()

	log.Debug("Neuron written to database", "uid", neuron.uid, "timepoint", neuron.timepoint)
}

// NeuronExists checks if a neuron exists by the given uid and timepoint
func (n *Neuroscan) NeuronExists(uid string, timepoint int) (bool, error) {
	db, err := n.ConnectDB()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var exists bool

	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM neurons WHERE uid = ? AND timepoint = ?)", uid, timepoint).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// UpdateNeuron updates the neuron record in the database by the given uid, for the provided params
func (n *Neuroscan) UpdateNeuron(uid string, timepoint int, embryonic bool, filename string, fileHash string, developmentalStage sql.NullInt64) error {
	db, err := n.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	exists, err := n.NeuronExists(uid, timepoint)

	if err != nil {
		return err
	}

	if !exists {
		return errors.New("neuron does not exist")
	}

	query := "UPDATE neurons SET "

	set := []string{}

	args := []interface{}{}

	if embryonic {
		set = append(set, "embryonic = ?")
		args = append(args, embryonic)
	}

	if filename != "" {
		set = append(set, "filename = ?")
		args = append(args, filename)
	}

	if fileHash != "" {
		set = append(set, "file_hash = ?")
		args = append(args, fileHash)
	}

	setString := strings.Join(set, ", ")

	query += setString

	query += " WHERE uid = ? AND timepoint = ?"

	// debug the query
	log.Debug("Query", "query", query)

	args = append(args, timepoint, uid)

	_, err = db.Exec(query, args...)
	if err != nil {
		return err
	}

	// get the neuron ID
	var neuronID int

	err = db.QueryRow("SELECT id FROM neurons WHERE uid = ? AND timepoint = ?", uid, timepoint).Scan(&neuronID)
	if err != nil {
		return err
	}

	// update the developmental stage
	_, err = db.Exec("UPDATE neurons__developmental_stages SET `developmental-stage_id` = ? WHERE neuron_id = ?", developmentalStage, neuronID)
	if err != nil {
		return err
	}

	return nil
}

// CreateNeuron creates a new neuron record in the database
func (n *Neuroscan) CreateNeuron(uid string, embryonic bool, filename string, fileHash string, developmentalStage sql.NullInt64, timepoint int) error {
	db, err := n.ConnectDB()
	if err != nil {
		return err
	}

	defer db.Close()

	exists, err := n.NeuronExists(uid, timepoint)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("neuron already exists")
	}

	result, err := db.Exec("INSERT INTO neurons (uid, embryonic, filename, file_hash, timepoint) VALUES (?, ?, ?, ?, ?)", uid, embryonic, filename, fileHash, timepoint)
	if err != nil {
		return err
	}

	neuronID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// insert the accompanying developmental stage join

	_, err = db.Exec("INSERT INTO neurons__developmental_stages (neuron_id, `developmental-stage_id`) VALUES (?, ?)", neuronID, developmentalStage)
	if err != nil {
		return err
	}

	return nil
}

// parseNeuron takes filename and returns a neuron object
func parseNeuron(n *Neuroscan, filePath string) (Neuron, error) {

	// filename is the last part of the path
	filename := filepath.Base(filePath)

	// check if we can get the file hash
	filehash, err := HashFile(filePath)
	if err != nil {
		log.Error("Error hashing file", "err", err)
		return Neuron{}, err
	}

	// parse the filename, it looks like SVV_RIAL.gltf, the UID is everything after the underscore and without the extension
	uid := BuildUID(filename)

	timepoint, err := GetTimepoint(filePath)
	if err != nil {
		log.Error("Error getting timepoint", "err", err)
		return Neuron{}, err
	}

	// get the developmental stage from the filename
	developmentalStage, err := GetDevStage(filePath)
	if err != nil {
		log.Error("Error getting developmental stage", "err", err)
		return Neuron{}, err
	}

	devStage, err := n.GetDevStageByUID(developmentalStage)

	if err != nil {
		log.Error("Error getting dev stage ID", "err", err)
		return Neuron{}, err
	}

	neuron := Neuron{
		uid:                uid,
		embryonic:          false,
		filename:           filename,
		fileHash:           filehash,
		developmentalStage: devStage.id,
		timepoint:          timepoint,
	}

	return neuron, nil
}

// ProcessNeuron processes a neuron file from the path, writing it to the database
func ProcessNeuron(n *Neuroscan, filePath string) {
	neuron, err := parseNeuron(n, filePath)
	if err != nil {
		log.Error("Error parsing neuron", "err", err)
		return
	}

	neuron.writeToDB(n)
}
