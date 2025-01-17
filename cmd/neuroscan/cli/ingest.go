package cli

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"sync/atomic"

	"neuroscan/internal/database"
	"neuroscan/internal/domain"
	"neuroscan/internal/logging"
	"neuroscan/internal/repository"
	"neuroscan/internal/service"
	"neuroscan/internal/toolshed"
)

type IngestCmd struct {
	DirPath string `required:"" help:"Path to the directory" short:"d"`
	Verbose bool `optional:"" help:"Enable verbose logging" short:"v"`
	SkipExisting bool `optional:"" help:"Skip existing files" short:"s"`
	ThreadCount int `optional:"" help:"Number of threads to use" short:"t"`
	ProcessTypes []string `optional:"" help:"Types of entities to process" short:"p"`
	Clean bool `optional:"" help:"Clean the database before ingesting" short:"c"`
}

type Ingestor struct {
	neurons      int64
	synapses     int64
	contacts     int64
	cphates      int64
	nerveRings   int64
	scales       int64
	skipExisting bool
	debug        bool
	clean        bool
	processTypes []string
	DevStages    []domain.DevelopmentalStage
	threadCount  int
}

type ingestChannels struct {
	neurons    chan string
	contacts   chan string
	synapses   chan string
	cphates    chan string
	nerveRings chan string
	scales     chan string
}

type ingestWaitGroups struct {
	neurons    sync.WaitGroup
	contacts   sync.WaitGroup
	synapses   sync.WaitGroup
	cphates    sync.WaitGroup
	nerveRings sync.WaitGroup
	scales     sync.WaitGroup
}

// createIngestChannels creates the ingest channels
func createIngestChannels() *ingestChannels {
	return &ingestChannels{
		neurons:    make(chan string, 10_000),
		contacts:   make(chan string, 100_000),
		synapses:   make(chan string, 100_000),
		cphates:    make(chan string, 20),
		nerveRings: make(chan string, 20),
		scales:     make(chan string, 20),
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
		scales:     sync.WaitGroup{},
	}
}

func (cmd *IngestCmd) Run(ctx *context.Context) error {
	logger := logging.NewLoggerFromEnv()
	cntx := logging.WithLogger(*ctx, logger)

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	db, err := database.NewFromEnv(cntx)
	if err != nil {
		logger.Fatal().Err(err).Msg("🤯 failed to connect to database")
		return err
	}

	defer db.Close(cntx)

	n := &Ingestor{
		neurons: 0,
		synapses: 0,
		contacts: 0,
		cphates: 0,
		nerveRings: 0,
		scales: 0,
		skipExisting: cmd.SkipExisting,
		debug:        cmd.Verbose,
		clean:        cmd.Clean,
		processTypes: cmd.ProcessTypes,
		threadCount:  cmd.ThreadCount,
	}

	// if processTypes is empty, set it to all valid process types
	if len(n.processTypes) == 0 {
		n.processTypes = []string{"neurons", "contacts", "synapses", "cphate", "nerveRing", "scale"}
	}

	channels := createIngestChannels()
	waitGroups := createIngestWaitGroups()

	maxRoutines := toolshed.MaxParallelism()

	neuronRepo := repository.NewPostgresNeuronRepository(db.Pool)
	neuronService := service.NewNeuronService(neuronRepo)

	contactRepo := repository.NewPostgresContactRepository(db.Pool)
	contactService := service.NewContactService(contactRepo)

	synapseRepo := repository.NewPostgresSynapseRepository(db.Pool)
	synapseService := service.NewSynapseService(synapseRepo)

	cphateRepo := repository.NewPostgresCphateRepository(db.Pool)
	cphateService := service.NewCphateService(cphateRepo)

	nerveRingRepo := repository.NewPostgresNerveRingRepository(db.Pool)
	nerveRingService := service.NewNerveRingService(nerveRingRepo)

	scaleRepo := repository.NewPostgresScaleRepository(db.Pool)
	scaleService := service.NewScaleService(scaleRepo)

	if n.clean {
		for _, processType := range n.processTypes {
			switch processType {
			case "neurons":
				err := neuronService.TruncateNeurons(cntx)
				if err != nil {
					logger.Error().Err(err).Msg("Error truncating neurons")
				}
			case "contacts":
				err := contactService.TruncateContacts(cntx)
				if err != nil {
					logger.Error().Err(err).Msg("Error truncating contacts")
				}
			case "synapses":
				err := synapseService.TruncateSynapses(cntx)
				if err != nil {
					logger.Error().Err(err).Msg("Error truncating synapses")
				}
			case "cphate":
				err := cphateService.TruncateCphates(cntx)
				if err != nil {
					logger.Error().Err(err).Msg("Error truncating cphates")
				}
			case "nerveRing":
				err := nerveRingService.TruncateNerveRings(cntx)
				if err != nil {
					logger.Error().Err(err).Msg("Error truncating nerveRings")
				}
			case "scale":
				err := scaleService.TruncateScales(cntx)
				if err != nil {
					logger.Error().Err(err).Msg("Error truncating scales")
				}
			}
		}
	}

	for w := 1; w <= maxRoutines; w++ {
		go func() {
			for neuronPath := range channels.neurons {
				neuron := domain.Neuron{}
				err := neuron.Parse(neuronPath)
				if err != nil {
					logger.Error().Err(err).Str("path", neuronPath).Msg("Error parsing neuron")
					waitGroups.neurons.Done()
					continue
				}

				success, err := neuronService.IngestNeuron(cntx, neuron, n.skipExisting, n.debug)
				if err != nil {
					logger.Error().Err(err).Str("path", neuronPath).Msg("Error ingesting neuron")
					waitGroups.neurons.Done()
					continue
				}

				if success {
					atomic.AddInt64(&n.neurons, 1)
				}

				waitGroups.neurons.Done()
			}

			for contactPath := range channels.contacts {
				contact := domain.Contact{}
				err := contact.Parse(contactPath)
				if err != nil {
					logger.Error().Err(err).Str("path", contactPath).Msg("Error parsing contact")
					waitGroups.contacts.Done()
					continue
				}

				success, err := contactService.IngestContact(cntx, contact, n.skipExisting, n.debug)
				if err != nil {
					logger.Error().Err(err).Str("path", contactPath).Msg("Error ingesting contact")
					waitGroups.contacts.Done()
					continue
				}

				if success {
					atomic.AddInt64(&n.contacts, 1)
				}

				waitGroups.contacts.Done()
			}

			for synapsePath := range channels.synapses {
				synapse := domain.Synapse{}
				err := synapse.Parse(synapsePath)
				if err != nil {
					logger.Error().Err(err).Str("path", synapsePath).Msg("Error parsing synapse")
					waitGroups.synapses.Done()
					continue
				}

				success, err := synapseService.IngestSynapse(cntx, synapse, n.skipExisting, n.debug)
				if err != nil {
					logger.Error().Err(err).Str("path", synapsePath).Msg("Error ingesting synapse")
					waitGroups.synapses.Done()
					continue
				}

				if success {
					atomic.AddInt64(&n.synapses, 1)
				}

				waitGroups.synapses.Done()
			}

			for cphateDir := range channels.cphates {
				cphate := domain.Cphate{}
				err := cphate.Parse(cphateDir)
				if err != nil {
					logger.Error().Err(err).Str("path", cphateDir).Msg("Error parsing cphate")
					waitGroups.cphates.Done()
					continue
				}

				success, err := cphateService.IngestCphate(cntx, cphate, n.skipExisting, n.debug)
				if err != nil {
					logger.Error().Err(err).Str("path", cphateDir).Msg("Error ingesting cphate")
					waitGroups.cphates.Done()
					continue
				}

				if success {
					atomic.AddInt64(&n.cphates, 1)
				}

				waitGroups.cphates.Done()
			}

			for nerveRingPath := range channels.nerveRings {
				nerveRing := domain.NerveRing{}
				err := nerveRing.Parse(nerveRingPath)
				if err != nil {
					logger.Error().Err(err).Str("path", nerveRingPath).Msg("Error parsing nerveRing")
					waitGroups.nerveRings.Done()
					continue
				}

				success, err := nerveRingService.IngestNerveRing(cntx, nerveRing, n.skipExisting, n.debug)
				if err != nil {
					logger.Error().Err(err).Str("path", nerveRingPath).Msg("Error ingesting nerveRing")
					waitGroups.nerveRings.Done()
					continue
				}

				if success {
					atomic.AddInt64(&n.nerveRings, 1)
				}

				waitGroups.nerveRings.Done()
			}

			for scalePath := range channels.scales {
				scale := domain.Scale{}
				err := scale.Parse(scalePath)
				if err != nil {
					logger.Error().Err(err).Str("path", scalePath).Msg("Error parsing scale")
					waitGroups.scales.Done()
					continue
				}

				success, err := scaleService.IngestScale(cntx, scale, n.skipExisting, n.debug)
				if err != nil {
					logger.Error().Err(err).Str("path", scalePath).Msg("Error ingesting scale")
					waitGroups.scales.Done()
					continue
				}

				if success {
					atomic.AddInt64(&n.scales, 1)
				}

				waitGroups.scales.Done()
			}
		}()
	}

	err = n.walkDirFolder(cntx, cmd.DirPath, channels, waitGroups)
	if err != nil {
		logger.Error().Err(err).Msg("Error processing entities")
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

	waitGroups.scales.Wait()
	close(channels.scales)

	logger.Info().Msg("Done processing entities")
	logger.Info().Int64("count", n.neurons).Msg("Neurons ingested")
	logger.Info().Int64("count", n.contacts).Msg("Contacts ingested")
	logger.Info().Int64("count", n.synapses).Msg("Synapses ingested")
	logger.Info().Int64("count", n.cphates).Msg("Cphates ingested")
	logger.Info().Int64("count", n.nerveRings).Msg("NerveRings ingested")
	logger.Info().Int64("count", n.scales).Msg("Scales ingested")

	return nil
}

func (n *Ingestor) walkDirFolder(ctx context.Context, path string, channels *ingestChannels, waitGroups *ingestWaitGroups) error {
	logger := logging.FromContext(ctx)

	logger.Info().Str("path", path).Msg("Walking directory")
	return filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			logger.Error().Err(err).Msg("Error walking directory")
			return err
		}

		// depending on the type of file, we want to process it differently
		currentEntity, err := toolshed.GetEntityType(path)

		if err != nil {
			logger.Error().Err(err).Str("path", path).Msg("Error getting entity type, skipping")
		}

		// we want to skip directories
		if d.IsDir() && currentEntity != "cphate" {
			return nil
		}

		if d.IsDir() && currentEntity == "cphate" {
			logger.Debug().Str("path", path).Msg("Adding cphate dir to channel")
			waitGroups.cphates.Add(1)
			channels.cphates <- path
		}

		// if it's not a valid extension, skip it
		if !toolshed.ValidExtension(path, []string{".gltf"}) {
			return nil
		}

		if !slices.Contains(n.processTypes, currentEntity) {
			return nil
		}

		// switch case to handle different entity types
		switch currentEntity {
		case "neurons":
			logger.Debug().Str("path", path).Msg("Adding neuron to channel")
			waitGroups.neurons.Add(1)
			channels.neurons <- path
		case "contacts":
			logger.Debug().Str("path", path).Msg("Adding contact to channel")
			waitGroups.contacts.Add(1)
			channels.contacts <- path
		case "synapses":
			logger.Debug().Str("path", path).Msg("Adding synapse to channel")
			waitGroups.synapses.Add(1)
			channels.synapses <- path
		case "nerveRing":
			logger.Debug().Str("path", path).Msg("Adding nerveRing to channel")
			waitGroups.nerveRings.Add(1)
			channels.nerveRings <- path
		case "cphate":
			logger.Debug().Str("path", path).Msg("Skipping cphate")
		case "scale":
			logger.Debug().Str("path", path).Msg("Adding scale to channel")
			waitGroups.scales.Add(1)
			channels.scales <- path
		default:
			logger.Error().Str("type", currentEntity).Msg("Unknown entity type")
		}

		logger.Debug().Str("path", path).Msg("Processing file")

		return nil
	})
}