-- +goose Up
-- +goose StatementBegin
create table if not exists events
(
    id          bigint PRIMARY KEY,
    title       text      NOT NULL,
    start_date  timestamp NOT NULL,
    end_date    timestamp NOT NULL,
    user_id     integer   NOT NULL
);
create sequence events_sequence start 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists events;
drop sequence events_sequence;
-- +goose StatementEnd
