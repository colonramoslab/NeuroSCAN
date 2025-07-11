-- +goose Up
-- +goose StatementBegin
-- CREATE EXTENSION IF NOT EXISTS "pgcrypto";
--
create table videos (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  ulid varchar(255) unique not null,
  status TEXT NOT NULL CHECK (status IN ('queued', 'processing', 'completed', 'failed')),
  error_message TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  started_at TIMESTAMPTZ,
  completed_at TIMESTAMPTZ
);

create index idx_videos_status on videos(status);
create index idx_videos_created_at on videos(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table videos;
-- +goose StatementEnd
