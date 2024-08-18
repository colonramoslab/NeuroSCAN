package main

import (
	"flag"
	"github.com/charmbracelet/log"
	"ingestion/neuroscan"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var processTypes arrayFlags

	dirPath := flag.String("dir", "", "Path to the directory")
	dbUrl := flag.String("db-url", "", "Database URL")
	debug := flag.Bool("debug", false, "Enable debug mode")
	skipExisting := flag.Bool("skip-existing", false, "Skip existing files")
	flag.Var(&processTypes, "types", "Types of entities to process (neurons, contacts, nerve_rings), defaults to all entities")
	flag.Parse()

	if *dirPath == "" || *dbUrl == "" {
		log.Fatal("Please provide a directory path and a database path")
	}

	log.SetLevel(log.InfoLevel)

	// create a new neuroscan object
	app := neuroscan.NewNeuroscan()

	if *debug {
		app.SetDebug(true)
		log.SetLevel(log.DebugLevel)
	}

	if *skipExisting {
		app.SetSkipExisting(true)
	}

	if len(processTypes) > 0 {
		log.Debug("Processing types: ", "types", processTypes)
		app.SetProcessTypes(processTypes)
	} else {
		log.Debug("Processing all types")
		app.SetDefaultProcessTypes()
	}

	// if we have a db url, set it
	if *dbUrl != "" {
		app.SetDBUrl(*dbUrl)
		app.SetDBType("postgres")
		app.BuildConnectionPool()
	}

	log.Debug("Database URL: ", "db-url", *dbUrl)
	log.Debug("Directory path: ", "dir", *dirPath)

	err := app.LoadDevStages()

	if err != nil {
		log.Fatal("Error loading developmental stages", "err", err)
	}

	// walk the directory
	app.ProcessEntities(*dirPath)

	app.CloseConnectionPool()

	log.Info("Done")
}
