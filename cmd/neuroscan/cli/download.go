package cli

import (
	"context"
	"runtime"

	"neuroscan/pkg/logging"
	"neuroscan/pkg/storage"
)

type DownloadCmd struct {
	Bucket			string `required:"" help:"S3 bucket to download files from." short:"b"`
	Destination string `optional:"" help:"Destination to save the downloaded files." short:"d"`
	Prefix      string `optional:"" help:"Prefix to search for in the S3 bucket." short:"p"`
	Key         string `optional:"" help:"Key to search for in the S3 bucket." short:"k"`
	Workers		  int    `optional:"" help:"Number of workers to use for downloading files." short:"w"`
}

func (cmd *DownloadCmd) Run(ctx *context.Context) error {
	logger := logging.FromContext(*ctx)

	logger.Info().Msg("Downloading files from S3 bucket")

	if cmd.Bucket == "" {
		logger.Error().Msg("No bucket provided, exiting")
		return nil
	}

	if cmd.Destination == "" {
		logger.Info().Msg("No destination provided, using current directory")
	}

	if cmd.Prefix == "" && cmd.Key == "" {
		logger.Error().Msg("No prefix or key provided, exiting")
		return nil
	}

	if cmd.Prefix != "" {
		logger.Info().Msgf("Searching for files with prefix: %s", cmd.Prefix)
	}

	if cmd.Key != "" {
		logger.Info().Msgf("Searching for file with key: %s", cmd.Key)
	}

	workers := runtime.NumCPU()

	if cmd.Workers > 0 {
		workers = cmd.Workers
	}

	if workers > 1 {
		logger.Info().Msgf("Using %d workers for downloading files", workers)
	}

	err := storage.DownloadBucketFile(cmd.Bucket, cmd.Destination, cmd.Key, cmd.Prefix, workers)

	if err != nil {
		logger.Error().Err(err).Msg("Error downloading files")
	}

	return nil
}
