-- +goose Up
-- +goose StatementBegin
create table contacts (
  id int generated always as identity primary key,
  ulid varchar(255) unique not null,
  uid varchar(255) not null,
  timepoint int not null,
  filename varchar(255) not null,
  color jsonb not null,
  unique (uid, timepoint)
);

create index idx_contacts_uid on contacts using gin (uid gin_trgm_ops);
create index idx_contacts_timepoint on contacts(timepoint);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table contacts;
-- +goose StatementEnd
