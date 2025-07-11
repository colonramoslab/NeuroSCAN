package transcode

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/h2non/filetype"
)

func ConvertWebmToMp4(ctx context.Context, webmPath string) error {
	// the mp4 file needs to be stored in the same directory as the webm file with the same name, but with the .mp4 extension
	mp4Filename := strings.Replace(webmPath, ".webm", ".mp4", 1)

	// if the mp4 file already exists, return an error
	if _, err := os.Stat(mp4Filename); err == nil {
		return fmt.Errorf("mp4 file already exists: %s", mp4Filename)
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", "-y", "-i", webmPath, "-c:v", "libx264", "-preset", "fast", "-movflags", "+faststart", mp4Filename)
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
