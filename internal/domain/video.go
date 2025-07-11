package domain

import (
	"time"

	"neuroscan/internal/toolshed"
)

type VideoStatus string

const (
	VideoStatusPending    VideoStatus = "pending"
	VideoStatusQueued     VideoStatus = "queued"
	VideoStatusProcessing VideoStatus = "processing"
	VideoStatusCompleted  VideoStatus = "completed"
	VideoStatusFailed     VideoStatus = "failed"
	VideoULIDPrefix                   = "vid"
)

type Video struct {
	ID           string      `json:"id"`
	ULID         string      `json:"uid"`
	Status       VideoStatus `json:"status"`
	ErrorMessage *string     `json:"error_message,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    *time.Time  `json:"updated_at"`
	CompletedAt  *time.Time  `json:"completed_at,omitempty"`
}

func (v *Video) New() error {
	ulid := toolshed.CreateULID(VideoULIDPrefix)

	v.ULID = ulid
	v.Status = VideoStatusPending

	return nil
}
