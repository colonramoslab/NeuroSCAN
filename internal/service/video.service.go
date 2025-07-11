package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"neuroscan/internal/domain"
	"neuroscan/internal/repository"
	"neuroscan/pkg/logging"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

type VideoService interface {
	GetVideoByUUID(ctx context.Context, uuid string) (domain.Video, error)
	CreateVideo(ctx context.Context, video domain.Video) (domain.Video, error)
	DeleteVideo(ctx context.Context, uuid string) error
	UpdateVideo(ctx context.Context, video domain.Video) (domain.Video, error)
	TruncateVideos(ctx context.Context) error
	Store(ctx context.Context, v domain.Video, data []byte) error
	Notify(ctx context.Context, v domain.Video) error
	TranscodeProcessing(ctx context.Context, uuid string) error
	TranscodeSuccess(ctx context.Context, uuid string) error
	TranscodeError(ctx context.Context, uuid string, err string) error
}

type videoService struct {
	repo repository.VideoRepository
}

func NewVideoRepository(repo repository.VideoRepository) VideoService {
	return &videoService{
		repo: repo,
	}
}

func (s *videoService) GetVideoByUUID(ctx context.Context, uuid string) (domain.Video, error) {
	return s.repo.GetVideoByUUID(ctx, uuid)
}

func (s *videoService) CreateVideo(ctx context.Context, video domain.Video) (domain.Video, error) {
	return s.repo.CreateVideo(ctx, video)
}

func (s *videoService) DeleteVideo(ctx context.Context, uuid string) error {
	return s.repo.DeleteVideo(ctx, uuid)
}

func (s *videoService) UpdateVideo(ctx context.Context, video domain.Video) (domain.Video, error) {
	return s.repo.UpdateVideo(ctx, video)
}

func (s *videoService) TruncateVideos(ctx context.Context) error {
	return s.repo.TruncateVideos(ctx)
}

func (s *videoService) TranscodeProcessing(ctx context.Context, uuid string) error {
	return s.repo.TranscodeProcessing(ctx, uuid)
}

func (s *videoService) TranscodeSuccess(ctx context.Context, uuid string) error {
	return s.repo.TranscodeSuccess(ctx, uuid)
}

func (s *videoService) TranscodeError(ctx context.Context, uuid string, err string) error {
	return s.repo.TranscodeError(ctx, uuid, err)
}

func (s *videoService) Store(ctx context.Context, v domain.Video, data []byte) error {
	logger := logging.NewLoggerFromEnv()

	err := godotenv.Load()
	if err != nil {
		logger.Info().Err(err).Msg("🤯 failed to load environment variables")
	}

	if v.ID == "" {
		return fmt.Errorf("video ID is required: %w", os.ErrInvalid)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	envDir := os.Getenv("VIDEO_STORAGE_PATH")

	if envDir == "" {
		return fmt.Errorf("VIDEO_STORAGE_PATH not set")
	}

	logger.Debug().Msgf("Storing video in %s", envDir)

	if err := os.MkdirAll(envDir, 0o755); err != nil {
		return fmt.Errorf("creating storage dir: %w", err)
	}

	filename := filepath.Join(homeDir, envDir, v.ID+".webm")

	logger.Debug().Msgf("Storing video as %s", filename)

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o644)
	if errors.Is(err, os.ErrExist) {
		return err
	}

	logger.Debug().Msgf("Opened file %s for writing", f.Name())

	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	logger.Debug().Msgf("Writing %d bytes to file %s", len(data), filename)
	if _, err = io.Copy(f, bytes.NewReader(data)); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	logger.Debug().Msgf("Syncing file %s to disk", filename)
	if err = f.Sync(); err != nil {
		return fmt.Errorf("fsync: %w", err)
	}

	logger.Debug().Msgf("File %s written and synced successfully", filename)
	v.Status = domain.VideoStatusQueued
	if _, err = s.UpdateVideo(ctx, v); err != nil {
		return fmt.Errorf("updating video status: %w", err)
	}

	logger.Debug().Msgf("Video %s stored successfully", v.ID)
	return nil
}

func (s *videoService) Notify(ctx context.Context, v domain.Video) error {
	logger := logging.NewLoggerFromEnv()

	err := godotenv.Load()
	if err != nil {
		logger.Info().Err(err).Msg("🤯 failed to load environment variables")
	}

	natsURL := os.Getenv("NATS_SERVER")
	if natsURL == "" {
		logger.Fatal().Msg("🤯 NATS_SERVER environment variable is not set")
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS server: %w", err)
	}

	defer nc.Close()

	nc.Publish("neuroscan.videos", []byte(v.ID))

	return nil
}
