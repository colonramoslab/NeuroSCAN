-- +goose Up
-- +goose StatementBegin
create table developmental_stages (
  id int generated always as identity primary key,
  uid varchar(255) not null,
  "begin" int not null,
  "end" int not null,
  "order" int not null,
  promoter_db boolean not null default false,
  timepoints int[] not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table developmental_stages;
-- +goose StatementEnd
