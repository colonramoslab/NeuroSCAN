-- +goose Up
-- +goose StatementBegin
create table promoters (
  id int generated always as identity primary key,
  uid varchar(255) unique not null,
  wormbase varchar(255),
  cellular_expression_pattern text,
  timepoint_start int,
  timepoint_end int,
  cells_by_lineaging varchar(255),
  expression_patterns text,
  information text,
  other_cells text,
  unique (uid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table promoters;
-- +goose StatementEnd
