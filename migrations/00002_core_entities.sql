-- +goose Up

-- Statuses (общий справочник статусов)
create table if not exists statuses (
    id serial primary key,
    name text not null unique
);

-- Пользователи
create table if not exists users (
    id serial primary key,
    name text not null,
    lastname text not null,
    status_id integer not null references statuses(id),
    rules text not null default '' -- роли/права в свободном формате (можно заменить на jsonb/массив позже)
);

-- Доски
create table if not exists boards (
    id serial primary key,
    name text not null unique
);

-- Тикеты
create table if not exists tickets (
    id serial primary key,
    name text not null,
    board_id integer not null references boards(id) on delete cascade,
    status_id integer not null references statuses(id),
    history jsonb not null default '[]'::jsonb,
    metadata jsonb not null default '{}'::jsonb,
    assignee_user_id integer references users(id) on delete set null
);

-- Карта переходов статусов (map: id_status -> next_status)
create table if not exists status_map (
    id_status integer not null references statuses(id) on delete cascade,
    next_status integer not null references statuses(id) on delete cascade,
    primary key (id_status, next_status),
    check (id_status <> next_status)
);

-- Индексы под типовые запросы
create index if not exists idx_users_status_id on users(status_id);
create index if not exists idx_tickets_board_id on tickets(board_id);
create index if not exists idx_tickets_status_id on tickets(status_id);
create index if not exists idx_tickets_assignee_user_id on tickets(assignee_user_id);
create index if not exists idx_status_map_id_status on status_map(id_status);

-- +goose Down

drop index if exists idx_status_map_id_status;
drop index if exists idx_tickets_assignee_user_id;
drop index if exists idx_tickets_status_id;
drop index if exists idx_tickets_board_id;
drop index if exists idx_users_status_id;

drop table if exists status_map;
drop table if exists tickets;
drop table if exists boards;
drop table if exists users;
drop table if exists statuses;

