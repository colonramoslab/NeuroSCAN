-- +goose Up
-- +goose StatementBegin
create table cphates (
  id int generated always as identity primary key,
  uid varchar(255) not null,
  timepoint int not null,
  structure jsonb not null,
  filename varchar(255) not null,
  unique (uid, timepoint)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table cphates;
-- +goose StatementEnd
