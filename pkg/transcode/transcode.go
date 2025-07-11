package transcode

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"neuroscan/pkg/logging"

	"github.com/h2non/filetype"
	"github.com/joho/godotenv"
)

func ConvertWebmToMp4(ctx context.Context, uuid string) error {
	logger := logging.NewLoggerFromEnv()

	err := godotenv.Load()
	if err != nil {
		logger.Info().Err(err).Msg("🤯 failed to load environment variables")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	videoDir := os.Getenv("VIDEO_STORAGE_PATH")

	if videoDir == "" {
		return fmt.Errorf("VIDEO_STORAGE_PATH environment variable is not set")
	}

	filename := filepath.Join(homeDir, videoDir, uuid+".webm")

	// if the mp4 file does not exist, return an error
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		return fmt.Errorf("webm file does not exist: %s", filename)
	}

	mp4Filename := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".mp4"

	cmd := exec.CommandContext(ctx, "ffmpeg", "-y", "-i", filename, "-c:v", "libx264", "-preset", "fast", "-movflags", "+faststart", mp4Filename)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg error: %w, details: %s", err, stderr.String())
	}

	return nil
}

func DetectFileType(buf []byte) (string, error) {
	kind, err := filetype.Match(buf)
	if err != nil {
		return "", fmt.Errorf("error detecting file type: %w", err)
	}
	if kind == filetype.Unknown {
		return "", fmt.Errorf("unknown file type")
	}
	return kind.MIME.Value, nil // e.g., "video/webm" or "video/mp4"
}
