-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_trgm;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop extension if exists pg_trgm;
-- +goose StatementEnd
