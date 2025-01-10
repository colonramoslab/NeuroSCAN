package toolshed

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"neuroscan/internal/gltf"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
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
			log.Debug().Int("timepoint", tp).Msg("Getting timepoint")
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
			log.Debug().Str("stage", part).Msg("Getting development stage")
			return part, nil
		}
	}

	// if no development stage is found, return an error
	return "", errors.New("development stage not found in path")
}

// GetEntityType returns the entity type based on the folder name in the filepath
func GetEntityType(filePath string) (string, error) {
	log.Debug().Str("path", filePath).Msg("Getting entity type")
	// the entity type is one of the following: neurons, synapses, contacts, cphate, nerveRing
	// it's based on the name of a folder somewhere in the middle of the path

	// split the path into parts
	parts := strings.Split(filePath, "/")

	// iterate over the parts
	for _, part := range parts {
		// if the part is one of the entity types, return it
		if part == "neurons" || part == "synapses" || part == "contacts" || part == "cphate" || part == "nerveRing" {
			log.Debug().Str("type", part).Msg("Getting entity type")
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
		log.Error().Err(err).Msg("Error getting file hash")
		return []NeuroscanFilepathData{}, err
	}

	timepoint, err := GetTimepoint(filePath)
	if err != nil {
		log.Error().Err(err).Msg("Error getting timepoint")
		return []NeuroscanFilepathData{}, err
	}

	devStageUID, err := GetDevStage(filePath)

	if err != nil {
		log.Error().Err(err).Msg("Error getting developmental stage")
		return []NeuroscanFilepathData{}, err
	}

	var parsedFiles []NeuroscanFilepathData
	// attempt to open and decode the gltf file
	doc, err := gltf.Open(filePath)
	if err != nil {
		log.Error().Err(err).Msg("Error opening gltf file")
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

func CreateDirectory(path string, permissions os.FileMode) error {
	// create directory if it doesn't exist at path
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, permissions)

		if err != nil {
			log.Error().Err(err).Msg("Error creating directory")
			return err
		}

		return nil
	}

	// if directory exists, return the path
	return nil
}

func BuildFilePath(dir string, fileName string, extension string) string {
	var path string

	// if dir is empty, use the current directory
	if dir == "" {
		dir = "."
	}

	// make sure dir does not have a trailing slash
	dir = filepath.Clean(dir)

	// if dir does not exist, create it
	err := CreateDirectory(dir, 0755)
	if err != nil {
		log.Error().Err(err).Msg("Error creating directory")
		return path
	}

	// if the fileName has a slash, remove it
	if filepath.Dir(fileName) != "." {
		fileName = filepath.Base(fileName)
	}

	// if the filename has an extension, remove it
	if filepath.Ext(fileName) != "" {
		fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	}

	// if fileName has no extension, create a manifest filename
	if filepath.Ext(fileName) == "" {
		fileName = fileName + extension
	}

	// create the file path
	path = filepath.Join(dir, fileName)

	return path
}

func CreateTempDirectory() string {
	// create a temp directory
	tempDirectory := os.TempDir() + "vsa-" + strconv.FormatInt(time.Now().Unix(), 10)
	CreateDirectory(tempDirectory, 0755)

	return tempDirectory
}

func GzipFile(input string, output string) error {
	// make sure the file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		return err
	}

	// open file
	file, err := os.Open(input)
	if err != nil {
		return err
	}
	defer file.Close()

	// create the directory if it doesn't exist
	CreateDirectory(filepath.Dir(output), 0755)

	// if output equals input, add .gz to the end
	if output == input {
		output = output + ".gz"
	}

	// create the gzip file
	gzipFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer gzipFile.Close()

	// create a gzip writer
	gzipWriter := gzip.NewWriter(gzipFile)
	defer gzipWriter.Close()

	// copy the file to the gzip writer
	_, err = io.Copy(gzipWriter, file)
	if err != nil {
		return err
	}

	_ = gzipWriter.Flush()

	return nil
}

func RemoveExtension(str string) string {
	return strings.Split(str, ".")[0]
}