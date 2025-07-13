package transcode

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"neuroscan/internal/cache"
	"neuroscan/internal/database"
	"neuroscan/internal/repository"
	"neuroscan/internal/service"
	"neuroscan/pkg/storage"

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
		logger.Info().Err(err).Msg("🤯 failed to load environment variables")
	}

	cntx := logging.WithLogger(ctxCancel, logger)

	natsURL := os.Getenv("NATS_SERVER")
	if natsURL == "" {
		logger.Fatal().Msg("🤯 NATS_SERVER environment variable is not set")
	}

	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		bucket = "neuroscan"
	}

	store, err := storage.NewStorage()
	if err != nil {
		logger.Fatal().Err(err).Msg("🤯 failed to create storage client")
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
	videoService := service.NewVideoService(videoRepo, *store, bucket)

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

	natsSubject := os.Getenv("NATS_VIDEO_SUBJECT")
	if natsSubject == "" {
		natsSubject = "neuroscan.videos"
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS server: %w", err)
	}
	// defer nc.Drain()

	logger.Info().Msg("Connected to NATS server")

	_, err = nc.Subscribe(natsSubject, func(m *nats.Msg) {
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

	err = cmd.convertWebmToMp4(ctx, videoService, uuid)
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

func (cmd *TranscodeCmd) convertWebmToMp4(ctx context.Context, videoService service.VideoService, uuid string) error {
	key := fmt.Sprintf("videos/%s.webm", uuid)
	remote := fmt.Sprintf("https://neuroscan-spaces.nyc3.digitaloceanspaces.com/%s", key)

	// ---------- ffmpeg setup ----------
	ff := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-i", remote, // input
		"-vf", "scale=-2:1080", // resize preserving aspect ratio
		"-r", "24", // 24 fps
		"-f", "mp4", // raw MP4 muxer to stdout
		"-movflags", "frag_keyframe+empty_moov+default_base_moof",
		"pipe:1", // write to stdout
	)

	// command := exec.CommandContext(ctx, "ffmpeg", "-y", "-i", remoteFile, "-vf", "scale=-2:1080", "-r", "24", destTemp.Name())
	var ffErr bytes.Buffer
	ff.Stderr = &ffErr

	// ---------- pipe ffmpeg → S3 ----------
	pr, pw := io.Pipe() // in-memory streaming pipe

	// 1. ffmpeg goroutine: write MP4 bytes to pw
	go func() {
		defer pw.Close()
		ff.Stdout = pw
		if err := ff.Run(); err != nil {
			pw.CloseWithError(
				fmt.Errorf("ffmpeg:%w – %s", err, ffErr.String()))
		}
	}()

	// 2. uploader goroutine (optional: could run inline)
	mp4Key := fmt.Sprintf("videos/%s.mp4", uuid)

	storageHandler := videoService.StorageHandle()
	if err := storageHandler.PutFile(videoService.BucketName(), mp4Key, pr); err != nil {
		return fmt.Errorf("upload mp4: %w", err)
	}

	return nil
}
