package service

import (
	"context"
	"os"
	"path/filepath"

	"neuroscan/internal/domain"
	"neuroscan/internal/repository"
	"neuroscan/pkg/logging"

	"github.com/joho/godotenv"
)

type VideoService interface {
	GetVideoByUUID(ctx context.Context, uuid string) (domain.Video, error)
	CreateVideo(ctx context.Context, video domain.Video) (domain.Video, error)
	DeleteVideo(ctx context.Context, uuid string) error
	UpdateVideo(ctx context.Context, video domain.Video) (domain.Video, error)
	TruncateVideos(ctx context.Context) error
	Store(ctx context.Context, v domain.Video, data []byte) error
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

func (s *videoService) Store(ctx context.Context, v domain.Video, data []byte) error {
	logger := logging.NewLoggerFromEnv()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("🤯 failed to load environment variables")
	}

	if v.ID == "" {
		return os.ErrInvalid
	}

	envDir := os.Getenv("VIDEO_STORAGE_PATH")

	if envDir == "" {
		return os.ErrInvalid
	}

	_, err = os.Stat(envDir)

	if os.IsNotExist(err) {
		err = os.MkdirAll(envDir, 0o755)
		if err != nil {
			return err
		}
	}

	filename := filepath.Join(envDir, v.ID+".webm")

	_, err = os.Stat(filename)
	if err == nil {
		return os.ErrExist
	}

	videoFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer videoFile.Close()

	_, err = videoFile.Write(data)
	if err != nil {
		return err
	}

	v.Status = domain.VideoStatusQueued

	_, err = s.UpdateVideo(ctx, v)
	if err != nil {
		return err
	}

	return nil
}
