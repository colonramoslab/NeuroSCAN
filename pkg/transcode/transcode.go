package transcode

import (
	"fmt"

	"github.com/h2non/filetype"
)

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
