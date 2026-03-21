-- +goose Up
create type "user_role" as enum (
    'user',
    'admin'
);

create type "social_network" as enum (
    'vk',
    'telegram',
    'github'
);

create type "edu_grade" as enum (
    'bachelor',
    'master',
    'specialist',
    'postgraduate'
);

create table "statuses" (
    id int primary key generated always as identity,
    content varchar (127) not null
);

create table "key_skills" (
    id int primary key generated always as identity,
    name varchar (31) not null
);

create table "soft_skills" (
    id int primary key generated always as identity,
    name varchar (31) not null
);

create table "users" (
    id uuid primary key default gen_random_uuid(),
    photo_s3_key varchar(255) not null,
    first_name varchar(31) not null,
    last_name varchar(31) not null,
    bio varchar(255) not null,
    birth_date date not null,
    status_id int references statuses(id),
    role user_role not null
);

create table "socials" (
    id int primary key generated always as identity,
    user_id uuid references users(id),
    social social_network not null,
    link varchar(127) not null,
    is_preferred boolean not null
);

-- +goose Down
