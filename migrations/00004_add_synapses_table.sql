-- +goose Up
-- +goose StatementBegin
create type synapse_type as enum ('chemical', 'electrical');

create table synapses (
  id int generated always as identity primary key,
  uid varchar(255) not null,
  timepoint int not null,
  synapse_type synapse_type,
  filename varchar(255) not null,
  color jsonb not null,
  unique (uid, timepoint)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table synapses;
-- +goose StatementEnd
