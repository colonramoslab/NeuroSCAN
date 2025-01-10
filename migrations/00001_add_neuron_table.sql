-- +goose Up
-- +goose StatementBegin
create table neurons (
  id int generated always as identity primary key,
  filename varchar(255) not null,
  timepoint int not null,
  uid varchar(255) not null,
  color jsonb not null,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table neurons;
-- +goose StatementEnd
