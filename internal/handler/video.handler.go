package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"neuroscan/internal/domain"
	"neuroscan/internal/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/echo/v4"
)

type VideoHandler struct {
	videoService service.VideoService
}

func NewVideoHandler(videoService service.VideoService) *VideoHandler {
	return &VideoHandler{
		videoService: videoService,
	}
}

func (h *VideoHandler) UploadWebm(c echo.Context) error {
	const max = int64(100 << 20)
	req := c.Request()
	req.Body = http.MaxBytesReader(c.Response(), req.Body, max)

	data, err := io.ReadAll(req.Body)
	if err != nil {
		return c.String(http.StatusBadRequest, "read body failed")
	}

	video := domain.Video{}
	if err := video.New(); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to create video object: %v", err))
	}
	newVideo, err := h.videoService.CreateVideo(req.Context(), video)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to create video in database: %v", err))
	}

	// 3. Fire-and-forget upload from RAM.
	go func(buf []byte, v domain.Video) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		if err := h.videoService.Store(ctx, v, bytes.NewReader(buf)); err != nil {
			_ = h.videoService.TranscodeError(ctx, v.ID, err.Error())
		} else {
			_ = h.videoService.Notify(ctx, v)
		}
	}(data, newVideo)

	// 4. Immediate response; body already consumed.
	return c.JSON(http.StatusAccepted, newVideo)
}

func (h *VideoHandler) UploadStatus(c echo.Context) error {
	uuid := c.Param("uuid")
	if uuid == "" {
		return c.String(http.StatusBadRequest, "invalid video ID")
	}

	video, err := h.videoService.GetVideoByUUID(c.Request().Context(), uuid)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to get video: %v", err))
	}

	return c.JSON(http.StatusOK, video)
}

func (h *VideoHandler) DownloadMP4(c echo.Context) error {
	filename := c.Param("filename")
	if filename == "" {
		return c.String(http.StatusBadRequest, "invalid video ID")
	}

	key := fmt.Sprintf("videos/%s", filename)

	rangeHeader := c.Request().Header.Get("Range")

	input := &s3.GetObjectInput{
		Bucket: aws.String(h.videoService.BucketName()),
		Key:    aws.String(key),
	}
	if rangeHeader != "" {
		input.Range = aws.String(rangeHeader)
	}

	obj, err := h.videoService.StorageHandle().Client.GetObject(input)
	if err != nil {
		// Translate common S3 errors into appropriate HTTP codes.
		if strings.Contains(err.Error(), "NotFound") {
			return echo.ErrNotFound
		}
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}
	defer obj.Body.Close()

	if obj.ContentType != nil {
		c.Response().Header().Set(echo.HeaderContentType, *obj.ContentType)
	}
	if obj.ContentLength != nil && rangeHeader == "" {
		c.Response().Header().Set(echo.HeaderContentLength, fmt.Sprint(*obj.ContentLength))
	}
	if obj.ETag != nil {
		c.Response().Header().Set("ETag", *obj.ETag)
	}
	if rangeHeader != "" && obj.ContentRange != nil {
		c.Response().Header().Set("Content-Range", *obj.ContentRange)
		c.Response().WriteHeader(http.StatusPartialContent)
	}

	_, err = io.Copy(c.Response(), obj.Body)
	return err
}
