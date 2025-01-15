-- +goose Up
-- +goose StatementBegin
create table contacts (
  id int generated always as identity primary key,
  uid varchar(255) not null,
  timepoint int not null,
  filename varchar(255) not null,
  color jsonb not null,
  unique (uid, timepoint)
);

create index contacts_timepoint_idx on contacts(timepoint);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table contacts;
-- +goose StatementEnd
