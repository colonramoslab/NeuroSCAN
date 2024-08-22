package neuroscan

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/charmbracelet/log"
	"ingestion/gltf"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Neuroscan struct {
	neurons      int64
	synapses     int64
	contacts     int64
	cphates      int64
	nerveRings   int64
	skipExisting bool
	debug        bool
	dbUrl        string
	dbType       string
	processTypes []string
	DevStages    []DevStage
	threadCount  int
	context      context.Context
	connPool     *pgxpool.Pool
}

type ingestChannels struct {
	neurons    chan string
	contacts   chan string
	synapses   chan string
	cphates    chan string
	nerveRings chan string
}

type ingestWaitGroups struct {
	neurons    sync.WaitGroup
	contacts   sync.WaitGroup
	synapses   sync.WaitGroup
	cphates    sync.WaitGroup
	nerveRings sync.WaitGroup
}

type FilepathData struct {
	uid                string
	filename           string
	filehash           string
	timepoint          int
	developmentalStage string
	color              Color
}

type Color [4]float64

// NewNeuroscan creates a new neuroscan object
func NewNeuroscan() *Neuroscan {
	return &Neuroscan{
		neurons:      0,
		synapses:     0,
		contacts:     0,
		cphates:      0,
		nerveRings:   0,
		skipExisting: false,
		dbType:       "",
		dbUrl:        "",
		processTypes: []string{},
		DevStages:    []DevStage{},
		threadCount:  1,
		debug:        false,
		context:      context.Background(),
	}
}

// SetThreadCount sets the thread count in the Neuroscan object
func (n *Neuroscan) SetThreadCount(threadCount int) {
	n.threadCount = threadCount
}

// GetThreadCount gets the thread count from the Neuroscan object
func (n *Neuroscan) GetThreadCount() int {
	return n.threadCount
}

// IncrementNeuron increments the neuron count in the Neuroscan object
func (n *Neuroscan) IncrementNeuron() {
	atomic.AddInt64(&n.neurons, 1)
}

// IncrementContact increments the contact count in the Neuroscan object
func (n *Neuroscan) IncrementContact() {
	atomic.AddInt64(&n.contacts, 1)
}

// IncrementSynapse increments the synapse count in the Neuroscan object
func (n *Neuroscan) IncrementSynapse() {
	atomic.AddInt64(&n.synapses, 1)
}

// IncrementCphate increments the cphate count in the Neuroscan object
func (n *Neuroscan) IncrementCphate() {
	atomic.AddInt64(&n.cphates, 1)
}

// IncrementNerveRing increments the nerveRing count in the Neuroscan object
func (n *Neuroscan) IncrementNerveRing() {
	atomic.AddInt64(&n.nerveRings, 1)
}

// SetDebug sets the debug flag in the Neuroscan object
func (n *Neuroscan) SetDebug(debug bool) {
	n.debug = debug
}

// SetDBType sets the database type in the Neuroscan object
func (n *Neuroscan) SetDBType(dbType string) {
	n.dbType = dbType
}

// GetDBType gets the database type from the Neuroscan object
func (n *Neuroscan) GetDBType() string {
	return n.dbType
}

// GetDBUrl gets the database URL from the Neuroscan object
func (n *Neuroscan) GetDBUrl() string {
	return n.dbUrl
}

// SetDBUrl sets the database URL in the Neuroscan object
func (n *Neuroscan) SetDBUrl(url string) {
	n.dbUrl = url
}

// SetSkipExisting sets the skip existing flag in the Neuroscan object
func (n *Neuroscan) SetSkipExisting(skipExisting bool) {
	n.skipExisting = skipExisting
}

// SetProcessTypes sets the process types in the Neuroscan object
func (n *Neuroscan) SetProcessTypes(processTypes []string) {
	n.processTypes = processTypes
}

// SetDefaultProcessTypes sets the default process types in the Neuroscan object
func (n *Neuroscan) SetDefaultProcessTypes() {
	n.processTypes = []string{"neurons", "contacts", "synapses", "cphate", "nerveRing"}
}

// BuildConnectionPool builds the connection pool on the Neuroscan object
func (n *Neuroscan) BuildConnectionPool() {

	// if we already have a connection pool, return
	if n.connPool != nil {
		return
	}

	connPool, err := pgxpool.New(n.context, n.dbUrl)
	if err != nil {
		log.Fatal("Error connecting to database", "err", err)
	}

	n.connPool = connPool
}

// CloseConnectionPool closes the connection pool on the Neuroscan object
func (n *Neuroscan) CloseConnectionPool() {
	n.connPool.Close()
}

// LoadDevStages loads the developmental stages from the database
func (n *Neuroscan) LoadDevStages() error {
	devStages, err := n.GetDevStagesAll()

	if err != nil {
		return err
	}

	n.DevStages = devStages

	return nil
}

// ConnectDB connects to the database
//func (n *Neuroscan) ConnectDB(ctx context.Context) (*pgx.Conn, error) {
//	//log.Debug("Connecting to database", "path", n.dbPath)
//	conn, err := pgx.Connect(ctx, n.dbUrl)
//	if err != nil {
//		return nil, err
//	}
//
//	return conn, nil
//}

//func (n *Neuroscan) ConnectDB(ctx context.Context) error {
//	db, err := sql.Open("sqlite3", "./foo.db")
//
//	if err != nil {
//		return err
//	}
//
//	n.conn = db
//	return nil
//}
//
//func (n *Neuroscan) DisconnectDB(ctx context.Context) error {
//	err := n.conn.Close()
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

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

// GetTimeNow get the time in the format "2006-01-02 15:04:05.000000-07"
func GetTimeNow() string {
	now := time.Now()
	timeFormat := "2006-01-02 15:04:05.000000-07"
	return now.Format(timeFormat)
}

// FilePathParse takes a filepath and returns the various metadata relating to the context of the file
func FilePathParse(filePath string) ([]FilepathData, error) {
	filename := filepath.Base(filePath)

	filehash, err := HashFile(filePath)

	if err != nil {
		log.Error("Error getting file hash", "err", err)
		return []FilepathData{}, err
	}

	timepoint, err := GetTimepoint(filePath)
	if err != nil {
		log.Error("Error getting timepoint", "err", err)
		return []FilepathData{}, err
	}

	devStageUID, err := GetDevStage(filePath)

	if err != nil {
		log.Error("Error getting developmental stage", "err", err)
		return []FilepathData{}, err
	}

	var parsedFiles []FilepathData
	// attempt to open and decode the gltf file
	doc, err := gltf.Open(filePath)
	if err != nil {
		log.Error("Error opening gltf file", "err", err)
		return []FilepathData{}, err
	}

	color := doc.Materials[0].PBRMetallicRoughness.BaseColorFactor

	for _, node := range doc.Nodes {
		uid := node.Name

		parsedFile := FilepathData{
			uid:                uid,
			filename:           filename,
			filehash:           filehash,
			timepoint:          timepoint,
			developmentalStage: devStageUID,
			color:              *color,
		}

		parsedFiles = append(parsedFiles, parsedFile)

	}

	return parsedFiles, nil
}

// EmptyStringToNil converts an empty string to a nil string
func EmptyStringToNil(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// IntToNil converts 0 to nil
func IntToNil(i *int) sql.NullInt64 {
	if i == nil || *i == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: int64(*i),
		Valid: true,
	}
}

// Int64ToNil converts 0 to nil
func Int64ToNil(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

// walkDirFolder walks a directory and processes the files for a specific entity type
func (n *Neuroscan) walkDirFolder(path string, channels *ingestChannels, waitGroups *ingestWaitGroups) error {
	log.Info("Walking directory", "path", path)
	return filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Error("Error walking directory", "err", err)
			return err
		}

		// depending on the type of file, we want to process it differently
		currentEntity, err := GetEntityType(path)

		if err != nil {
			log.Error("Error getting entity type", "err", err)
		}

		// we want to skip directories
		if d.IsDir() && currentEntity != "cphate" {
			return nil
		}

		if d.IsDir() && currentEntity == "cphate" {
			log.Debug("Adding cphate dir to channel", "path", path)
			waitGroups.cphates.Add(1)
			channels.cphates <- path
		}

		// if it's not a valid extension, skip it
		if !ValidExtension(path) {
			return nil
		}

		if !slices.Contains(n.processTypes, currentEntity) {
			return nil
		}

		// switch case to handle different entity types
		switch currentEntity {
		case "neurons":
			log.Debug("Adding neuron to channel", "path", path)
			waitGroups.neurons.Add(1)
			channels.neurons <- path
		case "contacts":
			log.Debug("Adding contact to channel", "path", path)
			waitGroups.contacts.Add(1)
			channels.contacts <- path
		case "synapses":
			log.Debug("Adding synapse to channel", "path", path)
			waitGroups.synapses.Add(1)
			channels.synapses <- path
		case "nerveRing":
			log.Debug("Adding nerveRing to channel", "path", path)
			waitGroups.nerveRings.Add(1)
			channels.nerveRings <- path
		case "cphate":
			log.Debug("Skipping cphate", "path", path)
		default:
			log.Error("Unknown entity type", "type", currentEntity)
		}

		log.Debug("Processing file", "path", path)

		return nil
	})
}

// createIngestChannels creates the ingest channels
func createIngestChannels() *ingestChannels {
	return &ingestChannels{
		neurons:    make(chan string, 10_000),
		contacts:   make(chan string, 100_000),
		synapses:   make(chan string, 100_000),
		cphates:    make(chan string, 20),
		nerveRings: make(chan string, 20),
	}
}

// createIngestWaitGroups creates the ingest wait groups
func createIngestWaitGroups() *ingestWaitGroups {
	return &ingestWaitGroups{
		neurons:    sync.WaitGroup{},
		contacts:   sync.WaitGroup{},
		synapses:   sync.WaitGroup{},
		cphates:    sync.WaitGroup{},
		nerveRings: sync.WaitGroup{},
	}
}

// maxParallelism returns the maximum parallelism
func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()

	if maxProcs < numCPU {
		return maxProcs
	}

	return numCPU
}

// ProcessEntities processes the entities in the directory in the proper order
func (n *Neuroscan) ProcessEntities(path string) {

	channels := createIngestChannels()
	waitGroups := createIngestWaitGroups()

	maxRoutines := maxParallelism()

	// start the worker pool, we don't do multiple workers right now because sqlite3 does not handle concurrent writes
	for w := 1; w <= maxRoutines; w++ {
		go func() {
			for neuronPath := range channels.neurons {
				ProcessNeuron(n, neuronPath)
				waitGroups.neurons.Done()
			}

			for contactPath := range channels.contacts {
				ProcessContact(n, contactPath)
				waitGroups.contacts.Done()
			}

			for synapsePath := range channels.synapses {
				ProcessSynapse(n, synapsePath)
				waitGroups.synapses.Done()
			}

			for cphateDir := range channels.cphates {
				ProcessCphate(n, cphateDir)
				waitGroups.cphates.Done()
			}

			for nerveRingPath := range channels.nerveRings {
				ProcessNerveRing(n, nerveRingPath)
				waitGroups.nerveRings.Done()
			}
		}()
	}

	err := n.walkDirFolder(path, channels, waitGroups)
	if err != nil {
		log.Error("Error processing entities", "err", err)
	}

	waitGroups.neurons.Wait()
	close(channels.neurons)

	waitGroups.contacts.Wait()
	close(channels.contacts)

	waitGroups.synapses.Wait()
	close(channels.synapses)

	waitGroups.cphates.Wait()
	close(channels.cphates)

	waitGroups.nerveRings.Wait()
	close(channels.nerveRings)

	log.Info("Done processing entities")
	log.Info("Neurons ingested: ", "count", n.neurons)
	log.Info("Contacts ingested: ", "count", n.contacts)
	log.Info("Synapses ingested: ", "count", n.synapses)
	log.Info("Cphates ingested: ", "count", n.cphates)
	log.Info("NerveRings ingested: ", "count", n.nerveRings)
}
