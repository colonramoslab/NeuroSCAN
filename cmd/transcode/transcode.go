package transcode

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"neuroscan/internal/cache"
	"neuroscan/internal/database"
	"neuroscan/internal/repository"
	"neuroscan/internal/service"
	"neuroscan/pkg/transcode"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"

	"neuroscan/pkg/logging"
)

type TranscodeCmd struct{}

func (cmd *TranscodeCmd) Run(ctx *context.Context) error {
	logger := logging.NewLoggerFromEnv()

	ctxCancel, cancel := context.WithCancel(*ctx)

	err := godotenv.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("🤯 failed to load environment variables")
	}

	cntx := logging.WithLogger(ctxCancel, logger)

	natsURL := os.Getenv("NATS_SERVER")
	if natsURL == "" {
		logger.Fatal().Msg("🤯 NATS_SERVER environment variable is not set")
	}

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
		err := cmd.videoListener(cntx, videoService)
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

func (cmd *TranscodeCmd) videoListener(ctx context.Context, videoService service.VideoService) error {
	logger := logging.FromContext(ctx)
	natsURL := os.Getenv("NATS_SERVER")
	if natsURL == "" {
		logger.Fatal().Msg("🤯 NATS_SERVER environment variable is not set")
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS server: %w", err)
	}
	// defer nc.Drain()

	logger.Info().Msg("Connected to NATS server")

	_, err = nc.Subscribe("neuroscan.videos", func(m *nats.Msg) {
		logger.Info().Msgf("Received message for video ID: %s", string(m.Data))
		err = cmd.transcodeVideo(ctx, videoService, string(m.Data))
		if err != nil {
			logger.Error().Err(err).Msg("Error transcoding video")
			return
		}
		logger.Info().Msgf("Successfully processed video ID: %s", string(m.Data))
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to subject: %w", err)
	}

	return nil
}

func (cmd *TranscodeCmd) transcodeVideo(ctx context.Context, videoService service.VideoService, uuid string) error {
	fmt.Printf("Processing video: %s\n", uuid)
	err := videoService.TranscodeProcessing(ctx, uuid)
	if err != nil {
		return fmt.Errorf("error setting video %s to processing: %w", uuid, err)
	}

	err = transcode.ConvertWebmToMp4(ctx, uuid)
	if err != nil {
		_ = videoService.TranscodeError(ctx, uuid, err.Error())
		return fmt.Errorf("error converting video %s: %w", uuid, err)
	}

	err = videoService.TranscodeSuccess(ctx, uuid)
	if err != nil {
		return fmt.Errorf("error setting video %s to success: %w", uuid, err)
	}

	return nil
}
