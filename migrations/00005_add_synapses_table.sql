-- +goose Up
-- +goose StatementBegin

create table synapses (
  id int generated always as identity primary key,
  uid varchar(255) not null,
  timepoint int not null,
  synapse_type varchar(255),
  filename varchar(255) not null,
  color jsonb not null,
  unique (uid, timepoint)
);

create index idx_synapses_uid on synapses using gin (uid gin_trgm_ops);
create index idx_synapses_timepoint on synapses(timepoint);
create index idx_synapses_synapse_type on synapses(synapse_type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table synapses;
-- +goose StatementEnd
