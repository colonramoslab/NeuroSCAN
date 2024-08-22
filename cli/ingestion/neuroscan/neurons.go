package neuroscan

import (
	"errors"
	"github.com/charmbracelet/log"
)

type Neuron struct {
	id        int
	uid       string
	embryonic bool
	filename  string
	timepoint int
	color     Color
}

// GetNeuron gets the neuron by UID and returns it, taking an optional timepoint and developmental stage
// the timepoint defaults to 0, and the developmental stage defaults to 0
func (n *Neuroscan) GetNeuron(uid string, timepoint int) (Neuron, error) {
	log.Debug("Getting neuron", "uid", uid, "timepoint", timepoint)
	var neuron Neuron

	query := "SELECT neurons.id, uid, embryonic, filename, timepoint, color FROM neurons WHERE uid = $1 AND timepoint = $2"

	args := []interface{}{uid, timepoint}

	err := n.connPool.QueryRow(n.context, query, args...).Scan(&neuron.id, &neuron.uid, &neuron.embryonic, &neuron.filename, &neuron.timepoint, &neuron.color)
	if err != nil {
		log.Error("Error getting neuron: ", "uid", uid, "timepoint", timepoint, "err", err)
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

	// if the neurons exists and we skip existing, return
	if n.skipExisting && exists {
		log.Debug("Neuron exists, skipping", "uid", neuron.uid)
		return
	}

	if exists {
		err = n.DeleteNeuron(neuron.uid, neuron.timepoint)
		if err != nil {
			log.Error("Error deleting existing neuron", "err", err)
			return
		}
	}

	err = n.CreateNeuron(neuron.uid, neuron.embryonic, neuron.filename, neuron.timepoint, neuron.color)
	if err != nil {
		log.Error("Error inserting new neuron record", "err", err)
		return
	}

	n.IncrementNeuron()

	log.Debug("Neuron written to database", "uid", neuron.uid, "timepoint", neuron.timepoint)
}

// NeuronExists checks if a neuron exists by the given uid and timepoint
func (n *Neuroscan) NeuronExists(uid string, timepoint int) (bool, error) {
	var exists bool

	err := n.connPool.QueryRow(n.context, "SELECT EXISTS(SELECT 1 FROM neurons WHERE uid = $1 AND timepoint = $2)", uid, timepoint).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// UpdateNeuron updates the neuron record in the database by the given uid, for the provided params
//func (n *Neuroscan) UpdateNeuron(uid string, timepoint int, embryonic bool, filename string, fileHash string, developmentalStage sql.NullInt64) error {
//	conn, err := n.ConnectDB(n.context)
//	if err != nil {
//		return err
//	}
//	defer conn.Close(n.context)
//
//	exists, err := n.NeuronExists(uid, timepoint)
//
//	if err != nil {
//		return err
//	}
//
//	if !exists {
//		return errors.New("neuron does not exist")
//	}
//
//	query := "UPDATE neurons SET "
//
//	set := []string{}
//
//	args := []interface{}{}
//
//	if embryonic {
//		set = append(set, "embryonic = ?")
//		args = append(args, embryonic)
//	}
//
//	if filename != "" {
//		set = append(set, "filename = ?")
//		args = append(args, filename)
//	}
//
//	if fileHash != "" {
//		set = append(set, "file_hash = ?")
//		args = append(args, fileHash)
//	}
//
//	if developmentalStage.Valid {
//		set = append(set, "developmental_stage = ?")
//		args = append(args, developmentalStage)
//	}
//
//	setString := strings.Join(set, ", ")
//
//	query += setString
//
//	query += " WHERE uid = ? AND timepoint = ?"
//
//	// debug the query
//	log.Debug("Query", "query", query)
//
//	args = append(args, timepoint, uid)
//
//	_, err = conn.Exec(n.context, query, args...)
//	if err != nil {
//		return err
//	}
//
//	// get the neuron ID
//	var neuronID int
//
//	err = conn.QueryRow(n.context, "SELECT id FROM neurons WHERE uid = ? AND timepoint = ?", uid, timepoint).Scan(&neuronID)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

// DeleteNeuron deletes a neuron record from the database by the given uid and timepoint
func (n *Neuroscan) DeleteNeuron(uid string, timepoint int) error {
	if n.connPool == nil {
		return errors.New("connection pool is nil")
	}

	_, err := n.connPool.Exec(n.context, "DELETE FROM neurons WHERE uid = $1 AND timepoint = $2", uid, timepoint)
	if err != nil {
		return err
	}

	return nil
}

// CreateNeuron creates a new neuron record in the database
func (n *Neuroscan) CreateNeuron(uid string, embryonic bool, filename string, timepoint int, color Color) error {
	exists, err := n.NeuronExists(uid, timepoint)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("neuron already exists")
	}

	_, err = n.connPool.Exec(n.context, "INSERT INTO neurons (uid, embryonic, filename, timepoint, color) VALUES ($1, $2, $3, $4, $5)", uid, embryonic, filename, timepoint, color)
	if err != nil {
		return err
	}

	return nil
}

// parseNeuron takes filename and returns a neuron object
func parseNeuron(n *Neuroscan, filePath string) (Neuron, error) {
	fileMetas, err := FilePathParse(filePath)

	if err != nil {
		log.Error("Error parsing file path", "err", err)
		return Neuron{}, err
	}

	fileMeta := fileMetas[0]

	neuron := Neuron{
		uid:       fileMeta.uid,
		embryonic: false,
		filename:  fileMeta.filename,
		timepoint: fileMeta.timepoint,
		color:     fileMeta.color,
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
