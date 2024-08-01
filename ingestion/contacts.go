package main

import (
	"errors"
	"path/filepath"
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
	fileHash  string
	timepoint int
}

// GetContact gets the contact by UID and returns it, taking an optional timepoint
func (n *Neuroscan) GetContact(uid string, timepoint int) (Contact, error) {
	log.Debug("Getting contact", "uid", uid, "timepoint", timepoint)
	var contact Contact

	db, err := n.ConnectDB()

	if err != nil {
		log.Error("Error connecting to database", "err", err)
		return contact, err
	}

	defer db.Close()

	err = db.QueryRow("SELECT id, uid, weight, neuron_a, neuron_b, filename, file_hash, timepoint FROM contacts WHERE uid = ? AND timepoint = ?", uid, timepoint).Scan(&contact.id, &contact.uid, &contact.weight, IntToNil(&contact.neuronA), IntToNil(&contact.neuronB), &contact.filename, &contact.fileHash, &contact.timepoint)
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

	db, err := n.ConnectDB()

	if err != nil {
		log.Error("Error connecting to database", "err", err)
		return
	}

	defer db.Close()

	if exists {
		log.Debug("Contact exists, updating record", "uid", contact.uid)
		err := n.UpdateContact(contact.uid, contact.timepoint, contact.weight, contact.neuronA, contact.neuronB, contact.filename, contact.fileHash)

		if err != nil {
			log.Error("Error updating contact record", "err", err)
			return
		}
	} else {
		log.Debug("Contact does not exist, inserting new record", "uid", contact.uid)
		err := n.CreateContact(contact.uid, contact.weight, contact.neuronA, contact.neuronB, contact.filename, contact.fileHash, contact.timepoint)
		if err != nil {
			log.Error("Error inserting new contact record", "err", err)
			return
		}
	}

	n.IncrementContact()

	log.Debug("Contact written to database", "uid", contact.uid, "timepoint", contact.timepoint)
}

// ContactExists checks if a contact exists by the given uid and timepoint
func (n *Neuroscan) ContactExists(uid string, timepoint int) (bool, error) {
	db, err := n.ConnectDB()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var exists bool

	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM contacts WHERE uid = ? AND timepoint = ?)", uid, timepoint).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// UpdateContact updates a contact by the given uid, for the provided params
func (n *Neuroscan) UpdateContact(uid string, timepoint int, weight int, neuronA int, neuronB int, filename string, fileHash string) error {
	db, err := n.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	exists, err := n.ContactExists(uid, timepoint)

	if err != nil {
		return err
	}

	if !exists {
		return errors.New("contact does not exist")
	}

	query := "UPDATE contacts SET "

	set := []string{}

	args := []interface{}{}

	if weight != 0 {
		set = append(set, "weight = ?")
		args = append(args, weight)
	}

	if neuronA != 0 {
		set = append(set, "neuron_a = ?")
		args = append(args, neuronA)
	}

	if neuronB != 0 {
		set = append(set, "neuron_b = ?")
		args = append(args, neuronB)
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

	args = append(args, timepoint, uid)

	_, err = db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

// CreateContact creates a new contact in the database
func (n *Neuroscan) CreateContact(uid string, weight int, neuronA int, neuronB int, filename string, fileHash string, timepoint int) error {
	db, err := n.ConnectDB()
	if err != nil {
		return err
	}

	defer db.Close()

	exists, err := n.ContactExists(uid, timepoint)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("contact already exists")
	}

	_, err = db.Exec("INSERT INTO contacts (uid, weight, neuron_a, neuron_b, filename, file_hash, timepoint) VALUES (?, ?, ?, ?, ?, ?, ?)", uid, weight, neuronA, neuronB, filename, fileHash, timepoint)
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
	filename := filepath.Base(filePath)

	filehash, err := HashFile(filePath)

	if err != nil {
		log.Error("Error getting file hash", "err", err)
		return Contact{}, err
	}

	uid := BuildUID(filename)

	timepoint, err := GetTimepoint(filePath)
	if err != nil {
		log.Error("Error getting timepoint", "err", err)
		return Contact{}, err
	}

	devStageUID, err := GetDevStage(filePath)

	if err != nil {
		log.Error("Error getting developmental stage", "err", err)
		return Contact{}, err
	}

	devStage, err := n.GetDevStageByUID(devStageUID)

	if err != nil {
		log.Error("Error getting developmental stage", "err", err)
		return Contact{}, err
	}

	contactNeuronA, contactNeuronB := getContactUIDNeurons(uid)

	neuronA, err := n.GetNeuron(contactNeuronA, timepoint, devStage.id)

	if err != nil {
		log.Error("Error getting neuron", "err", err)
		neuronA = Neuron{
			id: 0,
		}
	}

	neuronB, err := n.GetNeuron(contactNeuronB, timepoint, devStage.id)

	if err != nil {
		log.Error("Error getting neuron", "err", err)
		neuronB = Neuron{
			id: 0,
		}
	}

	contact := Contact{
		uid:       uid,
		weight:    0,
		neuronA:   neuronA.id,
		neuronB:   neuronB.id,
		filename:  filename,
		fileHash:  filehash,
		timepoint: timepoint,
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
