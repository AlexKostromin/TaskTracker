-- +goose Up
create table if not exists trackers (
    id serial primary key,
    name text not null,
    description text
);

-- +goose Down
drop table if exists trackers;

