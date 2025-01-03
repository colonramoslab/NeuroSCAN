package neuroscan

import (
	"errors"
	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)

type NerveRing struct {
	id        int
	uid       string
	timepoint int
	filename  string
	color     Color
}

// GetNerveRing gets the nerve ring by UID and returns it, taking an optional timepoint and developmental stage
func (n *Neuroscan) GetNerveRing(uid string, timepoint int) (NerveRing, error) {
	var nerveRing NerveRing

	err := n.connPool.QueryRow(n.context, "SELECT id, uid, timepoint, filename, color FROM nerve_rings WHERE uid = $1 AND timepoint = $2", uid, timepoint).Scan(&nerveRing.id, &nerveRing.uid, &nerveRing.timepoint, &nerveRing.filename, &nerveRing.color)

	if err != nil {
		return nerveRing, err
	}

	return nerveRing, nil
}

// writeToDB writes the nerve ring to the database
func (nerveRing NerveRing) writeToDB(n *Neuroscan) {
	exists, err := n.NerveRingExists(nerveRing.uid, nerveRing.timepoint)

	if err != nil {
		log.Error("Error checking if nerve ring exists", "err", err)
		return
	}

	// if the contacts exists and we skip existing, return
	if n.skipExisting && exists {
		log.Debug("Nerve ring exists, skipping", "uid", nerveRing.uid)
		return
	}

	if exists {
		err = n.DeleteNerveRing(nerveRing.uid, nerveRing.timepoint)
		if err != nil {
			log.Error("Error deleting existing nerve ring", "err", err)
		}
	}

	name := "Nerve Ring " + strconv.Itoa(nerveRing.timepoint)

	err = n.CreateNerveRing(name, nerveRing.timepoint, nerveRing.filename, nerveRing.color)
	if err != nil {
		log.Error("Error creating nerve ring", "err", err)
		return
	}

	n.IncrementNerveRing()

	log.Debug("Nerve ring written to database", "uid", nerveRing.uid, "timepoint", nerveRing.timepoint)
}

// NerveRingExists checks if a nerve ring exists by UID and timepoint
func (n *Neuroscan) NerveRingExists(uid string, timepoint int) (bool, error) {
	var count int

	err := n.connPool.QueryRow(n.context, "SELECT COUNT(*) FROM nerve_rings WHERE uid = $1 AND timepoint = $2", uid, timepoint).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateNerveRing updates a nerve ring in the database
//func (n *Neuroscan) UpdateNerveRing(uid string, developmentalStage int, timepoint int, filename string, fileHash string) error {
//	conn, err := n.ConnectDB(n.context)
//
//	if err != nil {
//		return err
//	}
//	defer conn.Close(n.context)
//
//	exists, err := n.NerveRingExists(uid, timepoint)
//
//	if err != nil {
//		return err
//	}
//
//	if !exists {
//		return errors.New("nerve ring does nto exist")
//	}
//
//	query := "UPDATE nerve_rings SET "
//
//	set := []string{}
//
//	args := []interface{}{}
//
//	if developmentalStage != 0 {
//		set = append(set, "developmental_stage = ?")
//		args = append(args, developmentalStage)
//	}
//
//	if filename != "" {
//		set = append(set, "filename = ?")
//		args = append(args, filename)
//	}
//
//	if fileHash != "" {
//		set = append(set, "filehash = ?")
//		args = append(args, fileHash)
//	}
//
//	setString := strings.Join(set, ", ")
//
//	query += setString
//
//	query += " WHERE uid = ? AND timepoint = ?"
//
//	args = append(args, uid, timepoint)
//
//	_, err = conn.Exec(n.context, query, args...)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

// DeleteNerveRing deletes a nerve ring from the database
func (n *Neuroscan) DeleteNerveRing(uid string, timepoint int) error {
	_, err := n.connPool.Exec(n.context, "DELETE FROM nerve_rings WHERE uid = $1 AND timepoint = $2", uid, timepoint)
	if err != nil {
		return err
	}

	return nil
}

// CreateNerveRing creates a new nerve ring in the database
func (n *Neuroscan) CreateNerveRing(uid string, timepoint int, filename string, color Color) error {
	exists, err := n.NerveRingExists(uid, timepoint)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("nerve ring already exists")
	}

	_, err = n.connPool.Exec(n.context, "INSERT INTO nerve_rings (uid, timepoint, filename, color) VALUES ($1, $2, $3, $4)", uid, timepoint, filename, color)
	if err != nil {
		return err
	}

	return nil
}

// parseNerveRing parses a nerve ring file path and returns a bervering object
func parseNerveRing(n *Neuroscan, filePath string) (NerveRing, error) {
	fileMetas, err := FilePathParse(filePath)

	if err != nil {
		log.Error("Failed to parse file meta", "error", err)
		return NerveRing{}, err
	}

	fileMeta := fileMetas[0]

	return NerveRing{
		uid:       fileMeta.uid,
		timepoint: fileMeta.timepoint,
		filename:  fileMeta.filename,
		color:     fileMeta.color,
	}, nil

}

// ProcessNerveRing processes all nerve rings in a directory
func ProcessNerveRing(n *Neuroscan, filePath string) {
	nerveRing, err := parseNerveRing(n, filePath)

	if err != nil {
		log.Error("Failed to parse nerve ring", "error", err)
		return
	}

	nerveRing.writeToDB(n)
}
