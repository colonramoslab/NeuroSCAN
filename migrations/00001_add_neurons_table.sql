-- +goose Up
-- +goose StatementBegin
create table neurons (
  id int generated always as identity primary key,
  uid varchar(255) not null,
  timepoint int not null,
  filename varchar(255) not null,
  color jsonb not null
);

create index neurons_timepoint_idx on neurons(timepoint);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table neurons;
-- +goose StatementEnd
