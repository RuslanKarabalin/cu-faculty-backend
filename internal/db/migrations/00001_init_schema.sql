-- +goose Up
create type role_type as enum (
    'user',
    'admin'
);

create type social_type as enum (
    'vk',
    'telegram',
    'github'
);

-- +goose Down
