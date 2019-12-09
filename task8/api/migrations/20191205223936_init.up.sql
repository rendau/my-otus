create table event
(
    id         bigserial   not null
        constraint event_pk primary key,
    owner      text        not null,
    title      text        not null,
    text       text        not null default '',
    start_time timestamptz not null,
    end_time   timestamptz not null
);