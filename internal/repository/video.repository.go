package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"neuroscan/internal/cache"
	"neuroscan/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VideoRepository interface {
	GetVideoByUUID(ctx context.Context, uuid string) (domain.Video, error)
	CreateVideo(ctx context.Context, video domain.Video) (domain.Video, error)
	DeleteVideo(ctx context.Context, uuid string) error
	UpdateVideo(ctx context.Context, video domain.Video) (domain.Video, error)
	TruncateVideos(ctx context.Context) error
	TranscodeProcessing(ctx context.Context, uuid string) error
	TranscodeSuccess(ctx context.Context, uuid string) error
	TranscodeError(ctx context.Context, uuid string, err string) error
}

type Video struct {
	ID           string         `db:"id"`
	ULID         string         `db:"uid"`
	Status       sql.NullString `db:"status"`
	ErrorMessage sql.NullString `db:"error_message"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
	CompletedAt  sql.NullTime   `db:"completed_at"`
}

func (v *Video) ToDomain() domain.Video {
	video := domain.Video{
		ID:        v.ID,
		ULID:      v.ULID,
		CreatedAt: v.CreatedAt,
	}

	if v.Status.Valid {
		video.Status = domain.VideoStatus(v.Status.String)
	}

	if v.ErrorMessage.Valid {
		video.ErrorMessage = &v.ErrorMessage.String
	}

	if v.UpdatedAt.Valid {
		video.UpdatedAt = &v.UpdatedAt.Time
	}

	if v.CompletedAt.Valid {
		video.CompletedAt = &v.CompletedAt.Time
	}

	return video
}

type PostgresVideoRepository struct {
	cache cache.Cache
	DB    *pgxpool.Pool
}

func NewPostgresVideoRepository(db *pgxpool.Pool, c cache.Cache) *PostgresVideoRepository {
	return &PostgresVideoRepository{
		cache: c,
		DB:    db,
	}
}

func (r *PostgresVideoRepository) GetVideoByUUID(ctx context.Context, uuid string) (domain.Video, error) {
	query := "SELECT id, ulid, status, error_message, created_at, updated_at, completed_at FROM videos WHERE id = $1"

	var video domain.Video
	err := r.DB.QueryRow(ctx, query, uuid).Scan(&video.ID, &video.ULID, &video.Status, &video.ErrorMessage, &video.CreatedAt, &video.UpdatedAt, &video.CompletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Video{}, nil
		}

		return domain.Video{}, err
	}

	return video, nil
}

func (r *PostgresVideoRepository) CreateVideo(ctx context.Context, video domain.Video) (domain.Video, error) {
	query := `
		INSERT INTO videos (ulid, status, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, ulid, status, error_message, created_at, updated_at, completed_at
	`

	var v Video

	err := r.DB.QueryRow(ctx, query, video.ULID, video.Status).Scan(&v.ID, &v.ULID, &v.Status, &v.ErrorMessage, &v.CreatedAt, &v.UpdatedAt, &v.CompletedAt)
	if err != nil {
		return domain.Video{}, err
	}

	return v.ToDomain(), nil
}

func (r *PostgresVideoRepository) DeleteVideo(ctx context.Context, uuid string) error {
	query := "DELETE FROM videos WHERE id = $1"

	_, err := r.DB.Exec(ctx, query, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresVideoRepository) UpdateVideo(ctx context.Context, video domain.Video) (domain.Video, error) {
	query := `
	UPDATE videos SET status = $1, error_message = $2, updated_at = NOW(), completed_at = $3 WHERE ulid = $4
	RETURNING id, ulid, status, error_message, created_at, updated_at, completed_at
	`

	var v Video

	err := r.DB.QueryRow(ctx, query, video.Status, video.ErrorMessage, video.CompletedAt, video.ULID).Scan(&v.ID, &v.ULID, &v.Status, &v.ErrorMessage, &v.CreatedAt, &v.UpdatedAt, &v.CompletedAt)
	if err != nil {
		return domain.Video{}, err
	}

	return v.ToDomain(), nil
}

func (r *PostgresVideoRepository) TruncateVideos(ctx context.Context) error {
	query := "TRUNCATE TABLE videos RESTART IDENTITY CASCADE"

	_, err := r.DB.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresVideoRepository) TranscodeProcessing(ctx context.Context, uuid string) error {
	query := `
		UPDATE videos SET status = $1, error_message = NULL, updated_at = NOW() WHERE id = $2
	`

	_, err := r.DB.Exec(ctx, query, domain.VideoStatusProcessing, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresVideoRepository) TranscodeSuccess(ctx context.Context, uuid string) error {
	query := `
		UPDATE videos SET status = $1, error_message = NULL, updated_at = NOW(), completed_at = NOW() WHERE id = $2
	`

	_, err := r.DB.Exec(ctx, query, domain.VideoStatusCompleted, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresVideoRepository) TranscodeError(ctx context.Context, uuid string, errMsg string) error {
	query := `
		UPDATE videos SET status = $1, error_message = $2, updated_at = NOW(), completed_at = NOW() WHERE id = $2
	`

	_, err := r.DB.Exec(ctx, query, domain.VideoStatusFailed, uuid, errMsg)
	if err != nil {
		return err
	}

	return nil
}
