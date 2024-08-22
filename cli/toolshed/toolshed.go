package toolshed

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/charmbracelet/log"
	"io"
	"neuroscan/ingestion/gltf"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type NeuroscanFilepathData struct {
	UID                *string
	Filename           *string
	Filehash           *string
	Timepoint          *int
	DevelopmentalStage *string
}

// ValidExtension checks if the file has a valid extension, currently that's just .gltf
func ValidExtension(fileName string) bool {
	return filepath.Ext(fileName) == ".gltf"
}

// HashFile returns the SHA256 hash of a file
func HashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GetTimepoint returns the timepoint based on the folder name in the filepath
func GetTimepoint(filePath string) (int, error) {
	// the timepoint is one of the following: 0, 5, 8, 23, 27, 36, or 48
	// it's based on the name of a folder somewhere in the middle of the path

	// split the path into parts
	parts := strings.Split(filePath, "/")

	// iterate over the parts
	for _, part := range parts {
		// if the part is a number, return it
		if tp, err := strconv.Atoi(part); err == nil {
			log.Debug("Getting timepoint", "timepoint", tp)
			return tp, nil
		}
	}

	// if no timepoint is found, return an error
	return 0, errors.New("timepoint not found in path")
}

// GetDevStage returns the development stage based on the folder name in the filepath
func GetDevStage(filePath string) (string, error) {
	// the development stage is one of the following: L1, L2, L3, L4, Adult
	// it's based on the name of a folder somewhere in the middle of the path

	// split the path into parts
	parts := strings.Split(filePath, "/")

	// iterate over the parts
	for _, part := range parts {
		// if the part is one of the development stages, return it
		if part == "L1" || part == "L2" || part == "L3" || part == "L4" || part == "Adult" {
			log.Debug("Getting development stage", "stage", part)
			return part, nil
		}
	}

	// if no development stage is found, return an error
	return "", errors.New("development stage not found in path")
}

// GetEntityType returns the entity type based on the folder name in the filepath
func GetEntityType(filePath string) (string, error) {
	log.Debug("Getting entity type", "path", filePath)
	// the entity type is one of the following: neurons, synapses, contacts, cphate, nerveRing
	// it's based on the name of a folder somewhere in the middle of the path

	// split the path into parts
	parts := strings.Split(filePath, "/")

	// iterate over the parts
	for _, part := range parts {
		// if the part is one of the entity types, return it
		if part == "neurons" || part == "synapses" || part == "contacts" || part == "cphate" || part == "nerveRing" {
			log.Debug("Getting entity type", "type", part)
			return part, nil
		}
	}

	// if no entity type is found, return an error
	return "", errors.New("entity type not found in path")
}

// CleanFilename cleans the filename by removing the path and extension
func CleanFilename(fileName string) string {
	return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
}

// BuildUID Build the neuron UId from the filename
func BuildUID(filename string) string {
	// parse the filename, it looks like SVV_RIAL.gltf, the UID is everything after the underscore and without the extension
	uid := CleanFilename(filename)
	// remove the _ and everything before it
	uid = strings.Split(uid, "_")[1]

	return uid
}

// FilePathParse takes a filepath and returns the various metadata relating to the context of the file
func FilePathParse(filePath string) ([]NeuroscanFilepathData, error) {
	filename := filepath.Base(filePath)

	filehash, err := HashFile(filePath)

	if err != nil {
		log.Error("Error getting file hash", "err", err)
		return []NeuroscanFilepathData{}, err
	}

	timepoint, err := GetTimepoint(filePath)
	if err != nil {
		log.Error("Error getting timepoint", "err", err)
		return []NeuroscanFilepathData{}, err
	}

	devStageUID, err := GetDevStage(filePath)

	if err != nil {
		log.Error("Error getting developmental stage", "err", err)
		return []NeuroscanFilepathData{}, err
	}

	var parsedFiles []NeuroscanFilepathData
	// attempt to open and decode the gltf file
	doc, err := gltf.Open(filePath)
	if err != nil {
		log.Error("Error opening gltf file", "err", err)
		return []NeuroscanFilepathData{}, err
	}

	for _, node := range doc.Nodes {
		uid := node.Name

		parsedFile := NeuroscanFilepathData{
			UID:                &uid,
			Filename:           &filename,
			Filehash:           &filehash,
			Timepoint:          &timepoint,
			DevelopmentalStage: &devStageUID,
		}

		parsedFiles = append(parsedFiles, parsedFile)

	}

	return parsedFiles, nil
}
