package neuroscan

import (
	"database/sql"
	"errors"
	"github.com/charmbracelet/log"
	"regexp"
	"strconv"
	"strings"
)

type Synapse struct {
	id          int
	uid         string
	synapseType string
	section     string
	position    string
	neuronSite  int
	neuronPre   sql.NullInt64
	postNeuron  int
	postNeurons []int
	timepoint   int
	filename    string
}

type SynapseData struct {
	neuronPre   string
	synapseType string
	section     string
	position    string
	neuronSite  int
	postNeuron  string
	postNeurons []string
}

type SynapsePosition struct {
	section  string
	position string
	site     int
}

// GetSynapse get the synapse by UID
func (n *Neuroscan) GetSynapse(uid string) (Synapse, error) {
	log.Debug("Getting synapse", "uid", uid)
	var synapse Synapse

	err := n.connPool.QueryRow(n.context, "SELECT id, uid, type, section, position, \"neuronSite\", \"neuronPre\", timepoint, filename FROM synapses WHERE uid = $1", uid).Scan(&synapse.id, &synapse.uid, &synapse.synapseType, &synapse.section, &synapse.position, &synapse.neuronSite, &synapse.neuronPre, &synapse.timepoint, &synapse.filename)
	if err != nil {
		return synapse, err
	}

	// get the post neurons associated with the synapse
	var postNeuronIDs []int
	err = n.connPool.QueryRow(n.context, "SELECT id FROM synapses__neuron_post WHERE synapse_id = $1", synapse.id).Scan(&postNeuronIDs)
	if err != nil {
		return synapse, nil
	}

	synapse.postNeurons = postNeuronIDs

	return synapse, nil
}

// writeToDB writes or updates the synapse to the database
func (synapse Synapse) writeToDB(n *Neuroscan) {
	exists, err := n.SynapseExists(synapse.uid, synapse.timepoint)

	if err != nil {
		log.Error("Error checking if synapse exists", "err", err)
	}

	// if the synapses exists and we skip existing, return
	if n.skipExisting && exists {
		log.Debug("Neuron exists, skipping", "uid", synapse.uid)
		return
	}

	if exists {
		err = n.DeleteSynapse(synapse.uid, synapse.timepoint)
		if err != nil {
			log.Error("Error deleting existing synapse", "err", err)
		}
	}

	log.Debug("Synapse does not exist, creating record", "uid", synapse.uid)
	err = n.CreateSynapse(synapse.uid, synapse.synapseType, synapse.section, synapse.position, synapse.neuronSite, synapse.neuronPre, synapse.postNeurons, synapse.timepoint, synapse.filename)

	if err != nil {
		log.Error("Error creating synapse record", "err", err)
	}

	n.IncrementSynapse()

	log.Info("Synapse created", "uid", synapse.uid)
}

// SynapseExists checks if the synapse exists
func (n *Neuroscan) SynapseExists(uid string, timepoint int) (bool, error) {
	var exists bool

	err := n.connPool.QueryRow(n.context, "SELECT EXISTS(SELECT 1 FROM synapses where uid = $1 AND timepoint = $2)", uid, timepoint).Scan(&exists)
	if err != nil {
		return false, nil
	}

	return exists, nil
}

// CreateSynapse creates a new synapse and it's related postneurons
func (n *Neuroscan) CreateSynapse(uid string, synapseType string, section string, position string, neuronSite int, neuronPre sql.NullInt64, postNeurons []int, timepoint int, filename string) error {
	exists, err := n.SynapseExists(uid, timepoint)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("synapse already exists")
	}

	_, err = n.connPool.Exec(n.context, "INSERT INTO synapses (uid, \"neuronPre\", type, section, position, \"neuronSite\", timepoint, filename) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", uid, neuronPre, synapseType, section, position, neuronSite, timepoint, filename)

	if err != nil {
		return err
	}

	// ge the synapse by UID
	synapse, err := n.GetSynapse(uid)

	if err != nil {
		return err
	}

	// now we need in insert the post neurons
	for _, postNeuron := range postNeurons {
		_, err = n.connPool.Exec(n.context, "INSERT INTO synapses__neuron_post (synapse_id, neuron_id) VALUES ($1, $2)", synapse.id, postNeuron)
		if err != nil {
			log.Error("Error inserting post neuron", "err", err)
		}
	}
	// now we need in insert the post neurons
	for _, postNeuron := range postNeurons {
		_, err = n.connPool.Exec(n.context, "INSERT INTO synapses__neuron_post (synapse_id, neuron_id) VALUES ($1, $2)", synapse.id, postNeuron)
		if err != nil {
			log.Error("Error inserting post neuron", "err", err)
		}
	}

	return nil
}

// DeleteSynapse deletes the synapse and it's related post neurons
func (n *Neuroscan) DeleteSynapse(uid string, timepoint int) error {
	_, err := n.connPool.Exec(n.context, "DELETE FROM synapses__neuron_post WHERE synapse_id = (SELECT id FROM synapses WHERE uid = $1 AND timepoint = $2)", uid, timepoint)

	if err != nil {
		return err
	}

	_, err = n.connPool.Exec(n.context, "DELETE FROM synapses__neuron_post WHERE synapse_id = (SELECT id FROM synapses WHERE uid = $1 AND timepoint = $2)", uid, timepoint)

	if err != nil {
		return err
	}

	_, err = n.connPool.Exec(n.context, "DELETE FROM synapses WHERE uid = $1 AND timepoint = $2", uid, timepoint)

	if err != nil {
		return err
	}

	return nil
}

// UpdateSynapse updates the synapse record
//func (n *Neuroscan) UpdateSynapse(uid string, synapseType string, section string, position string, neuronSite int, neuronPre sql.NullInt64, postNeurons []int, timepoint int, filename string, fileHash string) error {
//	conn, err := n.ConnectDB(n.context)
//
//	if err != nil {
//		return err
//	}
//	defer conn.Close(n.context)
//
//	exists, err := n.SynapseExists(uid, timepoint)
//
//	if err != nil {
//		return err
//	}
//
//	if !exists {
//		return errors.New("synapse does not exist")
//	}
//
//	query := "UPDATE synapses SET "
//
//	var set []string
//
//	var args []interface{}
//
//	if synapseType != "" {
//		set = append(set, "type = ?")
//		args = append(args, synapseType)
//	}
//
//	if section != "" {
//		set = append(set, "section = ?")
//		args = append(args, section)
//	}
//
//	if position != "" {
//		set = append(set, "position = ?")
//		args = append(args, position)
//	}
//
//	if neuronSite != 0 {
//		set = append(set, "neuron_site = ?")
//		args = append(args, neuronSite)
//	}
//
//	if neuronPre.Valid {
//		set = append(set, "neuron_pre = ?")
//		args = append(args, neuronPre)
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
//	// now clear the post neurons
//	_, err = conn.Exec(n.context, "DELETE FROM synapses__neuron_post WHERE synapse_id = $1", uid)
//
//	if err != nil {
//		return err
//	}
//
//	_, err = conn.Exec(n.context, "DELETE FROM synapse_post_neurons WHERE synapse_id = $1", uid)
//
//	if err != nil {
//		return err
//	}
//
//	// now we need in insert the post neurons
//	for _, postNeuron := range postNeurons {
//		_, err = conn.Exec(n.context, "INSERT INTO synapses__neuron_post (synapse_id, neuron_id) VALUES ($1, $2)", uid, postNeuron)
//		if err != nil {
//			log.Error("Error inserting post neuron", "err", err)
//		}
//
//		_, err = conn.Exec(n.context, "INSERT INTO synapse_post_neurons (synapse_id, neuron_id) VALUES ($1, $2)", uid, postNeuron)
//		if err != nil {
//			log.Error("Error inserting post neuron", "err", err)
//		}
//	}
//
//	return nil
//}

// buildSynapseTypeString takes a synapse UID and returns a valid type string
func buildSynapseTypeString(uid string) string {
	if strings.Contains(uid, "chemical") {
		return "chemical"
	}

	if strings.Contains(uid, "electrical") {
		return "electrical"
	}

	if strings.Contains(uid, "undefined") {
		return "undefined"
	}

	return ""
}

// buildSynapsePosition takes the position string and return a
func buildSynapsePosition(positionString string) SynapsePosition {
	// some examples are `A_post4`, `A_pre`, `B_pre`
	// the letter before the underscore is the section
	// if the digit at the end is present, that is the site
	// if present, pre or post is the position
	section := ""
	position := ""
	site := 0

	// if the string contains an _, then we try to parse a section
	if strings.Contains(positionString, "_") {
		section = strings.Split(positionString, "_")[0]

		// if section is a single letter between A - Z, then it is valid
		if len(section) == 1 {
			section = strings.ToUpper(section)
		}
	}

	// if the string contains a digit at the end, then we try to parse a site
	re := regexp.MustCompile(`\d+$`)
	match := re.FindString(positionString)

	if match != "" {
		siteInt, err := strconv.Atoi(match)

		if err != nil {
			log.Error("Error converting site to int", "err", err)
		} else {
			site = siteInt
		}
	}

	// if the string contains pre or post, then we try to parse a position
	if strings.Contains(positionString, "pre") {
		position = "pre"
	}

	if strings.Contains(positionString, "post") {
		position = "post"
	}

	return SynapsePosition{
		section:  section,
		position: position,
		site:     site,
	}
}

// getSynapseData takes a synapse uid string and parses out the component parts
// an example string is URYVRundefinedRMDDR&IL1VR~A_pre
// if we split it by the type, we should have URYVR and RMDDR&IL1VR~A_pre
// then we split the second part by the ~, so we should have RMDDR&IL1VR and A_pre
// the second part gets parsed for positioning, the first is our post neuron group
func getSynapseData(uid string) SynapseData {
	// first we will get the type, it's either chemical, electrical, or undefined. If not are present in the uid, it's nil
	synapseType := buildSynapseTypeString(uid)

	log.Debug("Synapse type", "type", synapseType)

	// if the Synapse type is not present, we cannot split the string by ""
	if synapseType == "" {
		log.Error("Synapse type not present in UID", "uid", uid)
	}

	neuronSections := strings.Split(uid, synapseType)
	log.Debug("Neuron sections", "sections", neuronSections)

	neuronPre := neuronSections[0]
	log.Debug("Neuron pre", "neuronPre", neuronPre)

	var positionString string

	if strings.Contains(neuronSections[1], "~") {
		positionStringParts := strings.Split(neuronSections[1], "~")
		log.Debug("Position string parts", "parts", positionStringParts)

		if len(positionStringParts) > 1 {
			positionString = positionStringParts[1]
		}
	}

	// if the position string is empty, just make it ""

	synapsePosition := buildSynapsePosition(positionString)
	log.Debug("Synapse position", "position", synapsePosition)

	postNeurons := strings.Split(strings.Split(neuronSections[1], "~")[0], "&")
	log.Debug("Post neurons", "postNeurons", postNeurons)

	return SynapseData{
		neuronPre:   neuronPre,
		synapseType: synapseType,
		section:     synapsePosition.section,
		position:    synapsePosition.position,
		neuronSite:  synapsePosition.site,
		postNeurons: postNeurons,
	}
}

// parseSynapse parses the synapse from the filename and returns a Synapse object
func parseSynapse(n *Neuroscan, filePath string) (Synapse, error) {
	fileMetas, err := FilePathParse(filePath)

	if err != nil {
		log.Error("Error parsing file path", "err", err)
		return Synapse{}, err
	}

	fileMeta := fileMetas[0]

	synapseData := getSynapseData(fileMeta.uid)

	//if we have a neuronPre from the synapse data, we can get the neuron ID
	neuronPreGet, err := n.GetNeuron(synapseData.neuronPre, fileMeta.timepoint)

	if err != nil {
		log.Error("Error getting neuron", "err", err)
		return Synapse{}, err
	}

	neuronPre := sql.NullInt64{
		Int64: int64(neuronPreGet.id),
		Valid: true,
	}

	var postNeuronIDs []int

	for _, postNeuron := range synapseData.postNeurons {
		neuron, err := n.GetNeuron(postNeuron, fileMeta.timepoint)

		if err != nil {
			log.Error("Error getting neuron", "err", err)
			return Synapse{}, err
		}

		postNeuronIDs = append(postNeuronIDs, neuron.id)
	}

	synapse := Synapse{
		uid:         fileMeta.uid,
		synapseType: synapseData.synapseType,
		section:     synapseData.section,
		position:    synapseData.position,
		neuronSite:  synapseData.neuronSite,
		neuronPre:   neuronPre,
		postNeurons: postNeuronIDs,
		timepoint:   fileMeta.timepoint,
		filename:    fileMeta.filename,
	}

	return synapse, nil
}

// ProcessSynapse processes a synapse from the path, writing it to the database
func ProcessSynapse(n *Neuroscan, filePath string) {
	synapse, err := parseSynapse(n, filePath)
	if err != nil {
		log.Error("Error parsing synapse", "err", err)
		return
	}

	synapse.writeToDB(n)
}
