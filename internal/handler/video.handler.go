package handler

import (
	"fmt"
	"io"
	"net/http"

	"neuroscan/internal/domain"
	"neuroscan/internal/service"

	"github.com/h2non/filetype"
	"github.com/labstack/echo/v4"
)

type VideoHandler struct {
	videoService service.VideoService
}

func NewVideoHandler(videoService service.VideoService) *VideoHandler {
	return &VideoHandler{videoService: videoService}
}

func (h *VideoHandler) UploadWebm(c echo.Context) error {
	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, "failed to read body")
	}

	// max size is 100MB
	if len(data) > 100*1024*1024 {
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

	video := domain.Video{}
	err = video.New()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to create video: %v", err))
	}

	newVideo, err := h.videoService.CreateVideo(c.Request().Context(), video)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to store video metadata: %v", err))
	}

	// store it in the filesystem
	err = h.videoService.Store(c.Request().Context(), newVideo, data)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to store video: %v", err))
	}

	// // Call conversion
	// mp4Bytes, err := ConvertWebmToMp4(c.Request().Context(), webmFile.Name())
	// if err != nil {
	// 	return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to convert to mp4: %v", err))
	// }

	// // Do something with mp4Bytes (store it, forward it, etc.)
	// return c.Blob(http.StatusOK, "video/mp4", mp4Bytes)
	//
	return c.JSON(http.StatusOK, newVideo)
}
