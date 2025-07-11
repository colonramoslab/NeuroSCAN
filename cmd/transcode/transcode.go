package transcode

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"neuroscan/internal/cache"
	"neuroscan/internal/database"
	"neuroscan/internal/repository"
	"neuroscan/internal/service"

	"github.com/joho/godotenv"

	"neuroscan/pkg/logging"
)

type TranscodeCmd struct {
	Interval time.Duration `optional:"" help:"Interval to check for new videos to transcode." short:"i" default:"1s"`
}

func (cmd *TranscodeCmd) Run(ctx *context.Context) error {
	logger := logging.NewLoggerFromEnv()

	ctxCancel, cancel := context.WithCancel(*ctx)

	err := godotenv.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("🤯 failed to load environment variables")
	}

	cntx := logging.WithLogger(ctxCancel, logger)

	db, err := database.NewFromEnv(cntx)
	if err != nil {
		logger.Fatal().Err(err).Msg("🤯 failed to connect to database")
		cancel()
		return err
	}
	defer db.Close(cntx)

	cache, err := cache.NewCache(cntx)
	if err != nil {
		logger.Fatal().Err(err).Msg("🤯 failed to connect to cache")
		cancel()
		return fmt.Errorf("failed to connect to cache: %w", err)
	}

	videoRepo := repository.NewPostgresVideoRepository(db.Pool, cache)
	videoService := service.NewVideoRepository(videoRepo)

	// we need to capture any interrupt signal to gracefully shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// WaitGroup to wait for goroutines to finish
	var wg sync.WaitGroup

	// Start your actual work in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := cmd.transcodeVideos(cntx, videoService)
		if err != nil {
			logger.Error().Err(err).Msg("Error during transcoding")
		}
	}()

	// Wait for a termination signal
	select {
	case sig := <-sigCh:
		fmt.Printf("\nReceived signal: %s. Initiating shutdown...\n", sig)
		cancel()
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("Shutdown complete.")

	return nil
}

func (cmd *TranscodeCmd) transcodeVideos(ctx context.Context, videoService service.VideoService) error {
	logger := logging.FromContext(ctx)
	ticker := time.NewTicker(cmd.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Graceful shutdown logic here
			logger.Info().Msg("Shutting down transcoding service...")
			return nil
		case t := <-ticker.C:
			logger.Info().Msgf("Checking for new videos to transcode at %v", t)
		}
	}
}
