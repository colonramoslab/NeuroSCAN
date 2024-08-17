package ingestion

import (
	"errors"
	"strings"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
)

type Contact struct {
	id        int
	uid       string
	weight    int
	neuronA   int `sql:",null"`
	neuronB   int `sql:",null"`
	filename  string
	timepoint int
}

// GetContact gets the contact by UID and returns it, taking an optional timepoint
func (n *Neuroscan) GetContact(uid string, timepoint int) (Contact, error) {
	log.Debug("Getting contact", "uid", uid, "timepoint", timepoint)
	var contact Contact

	err := n.connPool.QueryRow(n.context, "SELECT id, uid, weight, \"neuronA\", \"neuronB\", filename, timepoint FROM contacts WHERE uid = $1 AND timepoint = $2", uid, timepoint).Scan(&contact.id, &contact.uid, &contact.weight, IntToNil(&contact.neuronA), IntToNil(&contact.neuronB), &contact.filename, &contact.timepoint)
	if err != nil {
		return contact, err
	}

	return contact, nil
}

// writeToDB writes the contact to the database
func (contact Contact) writeToDB(n *Neuroscan) {
	exists, err := n.ContactExists(contact.uid, contact.timepoint)

	if err != nil {
		log.Error("Error checking if contact exists", "err", err)
		return
	}

	// if the contacts exists and we skip existing, return
	if n.skipExisting && exists {
		log.Debug("Neuron exists, skipping", "uid", contact.uid)
		return
	}

	if exists {
		// delete the contact
		err = n.DeleteContact(contact.uid, contact.timepoint)
		if err != nil {
			log.Error("Error deleting existing contact", "err", err)
			return
		}
	}

	log.Debug("Contact does not exist, inserting new record", "uid", contact.uid)
	err = n.CreateContact(contact.uid, contact.weight, contact.neuronA, contact.neuronB, contact.filename, contact.timepoint)
	if err != nil {
		log.Error("Error inserting new contact record", "err", err)
		return
	}

	n.IncrementContact()

	log.Debug("Contact written to database", "uid", contact.uid, "timepoint", contact.timepoint)
}

// ContactExists checks if a contact exists by the given uid and timepoint
func (n *Neuroscan) ContactExists(uid string, timepoint int) (bool, error) {
	var exists bool

	err := n.connPool.QueryRow(n.context, "SELECT EXISTS(SELECT 1 FROM contacts WHERE uid = $1 AND timepoint = $2)", uid, timepoint).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// UpdateContact updates a contact by the given uid, for the provided params
//func (n *Neuroscan) UpdateContact(uid string, timepoint int, weight int, neuronA int, neuronB int, filename string, fileHash string) error {
//	conn, err := n.ConnectDB(n.context)
//	if err != nil {
//		return err
//	}
//	defer conn.Close(n.context)
//
//	exists, err := n.ContactExists(uid, timepoint)
//
//	if err != nil {
//		return err
//	}
//
//	if !exists {
//		return errors.New("contact does not exist")
//	}
//
//	query := "UPDATE contacts SET "
//
//	set := []string{}
//
//	args := []interface{}{}
//
//	if weight != 0 {
//		set = append(set, "weight = ?")
//		args = append(args, weight)
//	}
//
//	if neuronA != 0 {
//		set = append(set, "neuronA = ?")
//		args = append(args, neuronA)
//	}
//
//	if neuronB != 0 {
//		set = append(set, "neuronB = ?")
//		args = append(args, neuronB)
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

// DeleteContact deletes a contact by the given uid and timepoint
func (n *Neuroscan) DeleteContact(uid string, timepoint int) error {
	_, err := n.connPool.Exec(n.context, "DELETE FROM contacts WHERE uid = $1 AND timepoint = $2", uid, timepoint)
	if err != nil {
		return err
	}

	return nil
}

// CreateContact creates a new contact in the database
func (n *Neuroscan) CreateContact(uid string, weight int, neuronA int, neuronB int, filename string, timepoint int) error {
	exists, err := n.ContactExists(uid, timepoint)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("contact already exists")
	}

	_, err = n.connPool.Exec(n.context, "INSERT INTO contacts (uid, weight, \"neuronA\", \"neuronB\", filename, timepoint) VALUES ($1, $2, $3, $4, $5, $6)", uid, weight, neuronA, neuronB, filename, timepoint)
	if err != nil {
		return err
	}

	return nil
}

// getContactUIDNeurons separates the contact filename using "by" in the string and returns the two neuron uid
func getContactUIDNeurons(filename string) (string, string) {
	// if the filename does not contain "by" then return empty strings
	if !strings.Contains(filename, "by") {
		return "", ""
	}

	neurons := strings.Split(filename, "by")
	return neurons[0], neurons[1]
}

// parseContact takes the file path and returns a contact object
func parseContact(n *Neuroscan, filePath string) (Contact, error) {
	fileMetas, err := FilePathParse(filePath)

	if err != nil {
		log.Error("Error parsing file path", "err", err)
		return Contact{}, err
	}

	fileMeta := fileMetas[0]

	contactNeuronA, contactNeuronB := getContactUIDNeurons(fileMeta.uid)

	neuronA, err := n.GetNeuron(contactNeuronA, fileMeta.timepoint)

	if err != nil {
		log.Error("Error getting neuron", "err", err)
		neuronA = Neuron{
			id: 0,
		}
	}

	neuronB, err := n.GetNeuron(contactNeuronB, fileMeta.timepoint)

	if err != nil {
		log.Error("Error getting neuron", "err", err)
		neuronB = Neuron{
			id: 0,
		}
	}

	contact := Contact{
		uid:       fileMeta.uid,
		weight:    0,
		neuronA:   neuronA.id,
		neuronB:   neuronB.id,
		filename:  fileMeta.filename,
		timepoint: fileMeta.timepoint,
	}

	return contact, nil
}

// ProcessContact processes a contact from the path, writing it to the database
func ProcessContact(n *Neuroscan, filePath string) {
	contact, err := parseContact(n, filePath)
	if err != nil {
		log.Error("Error parsing contact", "err", err)
		return
	}

	contact.writeToDB(n)
}
