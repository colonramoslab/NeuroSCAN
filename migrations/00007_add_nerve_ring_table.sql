-- +goose Up
-- +goose StatementBegin
create table nerve_rings (
  id int generated always as identity primary key,
  ulid varchar(255) unique not null,
  uid varchar(255) not null,
  timepoint int not null,
  filename varchar(255) not null,
  color jsonb not null,
  unique (uid, timepoint)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table nerve_rings;
-- +goose StatementEnd
