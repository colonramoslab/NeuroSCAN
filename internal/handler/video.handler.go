package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/h2non/filetype"
	"github.com/labstack/echo/v4"
)

func UploadWebm(c echo.Context) error {
	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, "failed to read body")
	}

	// max size is 20MB
	if len(data) > 20*1024*1024 {
		return c.String(http.StatusBadRequest, "file size exceeds limit")
	}

	kind, err := filetype.Match(data)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("error detecting file type: %v", err))
	}
	if kind == filetype.Unknown {
		return c.String(http.StatusBadRequest, "unknown file type")
	}
	if kind.MIME.Value != "video/webm" {
		return c.String(http.StatusBadRequest, "unsupported file type")
	}

	// Write to a temp file
	webmFile, err := os.CreateTemp("", "*.webm")
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to create temp file")
	}
	defer os.Remove(webmFile.Name())
	_, err = webmFile.Write(data)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to write webm")
	}
	webmFile.Close()

	// Call conversion
	mp4Bytes, err := ConvertWebmToMp4(c.Request().Context(), webmFile.Name())
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to convert to mp4: %v", err))
	}

	// Do something with mp4Bytes (store it, forward it, etc.)
	return c.Blob(http.StatusOK, "video/mp4", mp4Bytes)
}

func ConvertWebmToMp4(ctx context.Context, webmPath string) ([]byte, error) {
	mp4File, err := os.CreateTemp("", "*.mp4")
	if err != nil {
		return nil, fmt.Errorf("failed to create mp4 temp file: %w", err)
	}
	defer os.Remove(mp4File.Name())
	mp4File.Close()

	cmd := exec.CommandContext(ctx, "ffmpeg", "-y", "-i", webmPath, "-c:v", "libx264", "-preset", "veryfast", "-crf", "25", "-movflags", "+faststart", mp4File.Name())
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg error: %w, details: %s", err, stderr.String())
	}

	// Read the mp4 result back into memory
	return os.ReadFile(mp4File.Name())
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
