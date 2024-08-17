package neuroscan

import (
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"neuroscan/cmd/ingestion"
	"neuroscan/cmd/upload"
)

const (
	possibleArgs = "serve ingest upload"
)

func main() {
	args := os.Args
	if len(args) != 2 || !strings.Contains(possibleArgs, args[1]) {
		log.Fatal("Invalid arguments. Possible arguments: " + possibleArgs)
		os.Exit(1)
	}

	switch args[1] {
	case "ingest":
		ingestion.Run()
	case "upload":
		upload.Run()
	}
}
