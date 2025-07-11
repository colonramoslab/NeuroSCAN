package main

import (
	"context"
	"os"

	"neuroscan/cmd/ingest"
	"neuroscan/cmd/transcode"
	"neuroscan/cmd/web"
	"neuroscan/pkg/logging"

	"github.com/alecthomas/kong"
)

type Cli struct {
	Web       web.WebCmd             `cmd:"" help:"Start the web server."`
	Ingest    ingest.IngestCmd       `cmd:"" help:"Ingest files into the database."`
	Transcode transcode.TranscodeCmd `cmd:"" help:"Listen for videos and transcode."`
}

func main() {
	// Display help if no args are provided instead of an error message
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}

	app := Cli{}

	logger := logging.NewLoggerFromEnv()
	cntx := context.Background()

	ctx := kong.Parse(&app,
		kong.Name("neuroscan"),
		kong.Description("NeuroSCAN CLI"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)

	logging.WithLogger(cntx, logger)

	err := ctx.Run(&cntx)
	ctx.FatalIfErrorf(err)
}
