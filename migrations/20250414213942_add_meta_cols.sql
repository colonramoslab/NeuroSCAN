-- +goose Up
-- +goose StatementBegin
alter table neurons
add column volume double precision,
add column surface_area double precision;

alter table contacts
add column surface_area double precision;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
alter table neurons
drop column volume,
drop column surface_area;

alter table contacts
drop column surface_area;

-- +goose StatementEnd
