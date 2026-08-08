-- +goose Up
alter table "users" add column deleted_at timestamptz;

-- +goose Down
alter table "users" drop column if exists deleted_at;
