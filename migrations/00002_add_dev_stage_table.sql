-- +goose Up
-- +goose StatementBegin
create table developmental_stages (
  id int generated always as identity primary key,
  ulid varchar(255) unique not null,
  uid varchar(255) unique not null,
  "begin" int not null,
  "end" int not null,
  "order" int not null,
  promoter_db boolean,
  timepoints int[] not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table developmental_stages;
-- +goose StatementEnd
