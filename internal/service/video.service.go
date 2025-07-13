package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"neuroscan/internal/domain"
	"neuroscan/internal/repository"
	"neuroscan/pkg/logging"
	"neuroscan/pkg/storage"

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
	GetVideosOlderThan(ctx context.Context, cutoffTime time.Time) ([]domain.Video, error)
	StorageHandle() storage.Storage
	BucketName() string
}

type videoService struct {
	repo    repository.VideoRepository
	storage storage.Storage
	bucket  string
}

func NewVideoService(repo repository.VideoRepository, storage storage.Storage, bucket string) VideoService {
	return &videoService{
		repo:    repo,
		storage: storage,
		bucket:  bucket,
	}
}

func (s *videoService) StorageHandle() storage.Storage { return s.storage }
func (s *videoService) BucketName() string             { return s.bucket }

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

func (s *videoService) GetVideosOlderThan(ctx context.Context, cutoffTime time.Time) ([]domain.Video, error) {
	return s.repo.GetVideosOlderThan(ctx, cutoffTime)
}

func (s *videoService) Store(ctx context.Context, v domain.Video, data []byte) error {
	logger := logging.NewLoggerFromEnv()

	if v.ID == "" {
		return fmt.Errorf("video ID is required: %w", os.ErrInvalid)
	}

	key := fmt.Sprintf("videos/%s.webm", v.ID)

	err := s.storage.PutFile(s.bucket, key, data)
	if err != nil {
		return fmt.Errorf("uploading file to storage: %w", err)
	}

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
