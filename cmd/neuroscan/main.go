package main

import (
	"context"
	"os"

	"neuroscan/cmd/neuroscan/cli"
	"neuroscan/pkg/logging"

	"github.com/alecthomas/kong"
)

type Cli struct {
	Download cli.DownloadCmd `cmd:"" help:"Download files from the specified S3 bucket."`
	Web      cli.WebCmd      `cmd:"" help:"Start the web server."`
	Ingest   cli.IngestCmd   `cmd:"" help:"Ingest files into the database."`
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