package main

import (
	"flag"
	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dirPath := flag.String("dir", "", "Path to the directory")
	dbPath := flag.String("db", "", "Path to the database")
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *dirPath == "" || *dbPath == "" {
		log.Fatal("Please provide a directory path and a database path")
	}

	log.SetLevel(log.InfoLevel)

	// create a new neuroscan object
	neuroscan := NewNeuroscan()

	if *debug {
		neuroscan.SetDebug(true)
		log.SetLevel(log.DebugLevel)
	}

	// set the DB path
	neuroscan.SetDBPath(*dbPath)

	log.Debug("Database path: ", "db", *dbPath)
	log.Debug("Directory path: ", "dir", *dirPath)

	// walk the directory
	neuroscan.ProcessEntities(*dirPath)

	log.Info("Done")
}
