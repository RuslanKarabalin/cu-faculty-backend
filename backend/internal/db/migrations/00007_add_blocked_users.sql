-- +goose Up
create table "blocked_users" (
    user_id uuid references users(id) on delete cascade
    , blocked_user_id uuid references users(id) on delete cascade
    , primary key (user_id, blocked_user_id)
    , check (user_id <> blocked_user_id)
);

-- +goose Down
drop table if exists "blocked_users";
