package cleanup

import (
	"context"
	"fmt"
	"os"
	"time"

	"neuroscan/internal/cache"
	"neuroscan/internal/database"
	"neuroscan/internal/repository"
	"neuroscan/internal/service"
	"neuroscan/pkg/logging"
	"neuroscan/pkg/storage"

	"github.com/joho/godotenv"
)

type CleanupCmd struct{}

func (cmd *CleanupCmd) Run(ctx *context.Context) error {
	logger := logging.NewLoggerFromEnv()

	err := godotenv.Load()
	if err != nil {
		logger.Info().Err(err).Msg("🤯 failed to load environment variables")
	}

	cntx := logging.WithLogger(*ctx, logger)

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
		return err
	}
	defer db.Close(cntx)

	cache, err := cache.NewCache(cntx)
	if err != nil {
		logger.Fatal().Err(err).Msg("🤯 failed to connect to cache")
		return fmt.Errorf("failed to connect to cache: %w", err)
	}

	videoRepo := repository.NewPostgresVideoRepository(db.Pool, cache)
	videoService := service.NewVideoService(videoRepo, *store, bucket)

	return cmd.cleanupOldVideos(cntx, videoService)
}

func (cmd *CleanupCmd) cleanupOldVideos(ctx context.Context, videoService service.VideoService) error {
	logger := logging.FromContext(ctx)

	cutoffTime := time.Now().Add(-4 * time.Hour)
	logger.Info().Msgf("Starting cleanup of videos older than %s", cutoffTime.Format(time.RFC3339))

	oldVideos, err := videoService.GetVideosOlderThan(ctx, cutoffTime)
	if err != nil {
		return fmt.Errorf("failed to get old videos: %w", err)
	}

	logger.Info().Msgf("Found %d videos to clean up", len(oldVideos))

	for _, video := range oldVideos {
		logger.Info().Msgf("Cleaning up video: %s", video.ID)

		// Delete from storage first
		webmKey := fmt.Sprintf("videos/%s.webm", video.ID)
		mp4Key := fmt.Sprintf("videos/%s.mp4", video.ID)

		client := videoService.StorageHandle()

		// Try to delete both formats, don't fail if one doesn't exist
		if err := client.DeleteFile(videoService.BucketName(), webmKey); err != nil {
			logger.Warn().Err(err).Msgf("Failed to delete webm file for video %s", video.ID)
		}

		if err := client.DeleteFile(videoService.BucketName(), mp4Key); err != nil {
			logger.Warn().Err(err).Msgf("Failed to delete mp4 file for video %s", video.ID)
		}

		// Delete from database
		if err := videoService.DeleteVideo(ctx, video.ID); err != nil {
			logger.Error().Err(err).Msgf("Failed to delete video %s from database", video.ID)
			continue
		}

		logger.Info().Msgf("Successfully cleaned up video: %s", video.ID)
	}

	logger.Info().Msgf("Cleanup completed. Processed %d videos", len(oldVideos))
	return nil
}
