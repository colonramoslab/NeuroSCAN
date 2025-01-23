-- +goose Up
-- +goose StatementBegin
create table developmental_stages (
  id int generated always as identity primary key,
  ulid varchar(255) unique not null,
  uid varchar(255) unique not null,
  "begin" int not null,
  "end" int not null,
  "order" int not null,
  promoter_db boolean not null default false,
  timepoints int[] not null
);

insert into developmental_stages (uid, "begin", "end", "order", promoter_db, timepoints)
values
  ('L1', 0, 16, 2, false, ARRAY[0,5,8,16]),
  ('L2', 16, 24, 3, false, ARRAY[23]),
  ('L3', 24, 33, 4, false, ARRAY[27]),
  ('L4', 33, 45, 5, false, ARRAY[36]),
  ('Adult', 45, 60, 6, false, ARRAY[45,48]),
  ('2 cell', 50, 345, 1, true, ARRAY[50]),
  ('Bean', 345, 390, 2, true, ARRAY[345]),
  ('Comma', 390, 415, 3, true, ARRAY[390]),
  ('1.5 fold', 415, 435, 4, true, ARRAY[415]),
  ('Twitching', 435, 455, 5, true, ARRAY[435]),
  ('2 fold', 455, 550, 6, true, ARRAY[455]),
  ('3 fold', 550, 840, 7, true, ARRAY[550]);
  -- ('Hatching', 840, 900, 8, true, ARRAY[840]);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table developmental_stages;
-- +goose StatementEnd
