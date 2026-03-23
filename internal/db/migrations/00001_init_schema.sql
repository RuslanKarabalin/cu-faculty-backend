-- +goose Up
create type "user_role" as enum (
    'user'
    , 'admin'
);

create type "social_network" as enum (
    'vk'
    , 'telegram'
    , 'github'
);

create type "edu_grade" as enum (
    'bachelor'
    , 'master'
    , 'specialist'
    , 'postgraduate'
);

create table "statuses" (
    id int primary key generated always as identity
    , content varchar (127) not null
);

create table "key_skills" (
    id int primary key generated always as identity
    , name varchar (31) not null
);

create table "soft_skills" (
    id int primary key generated always as identity
    , name varchar (31) not null
);

create table "users" (
    id uuid primary key default gen_random_uuid()
    , photo_s3_key varchar(255) not null
    , first_name varchar(31) not null
    , last_name varchar(31) not null
    , bio varchar(255) not null
    , birth_date date not null
    , status_id int references statuses(id)
    , role user_role not null
);

create table "socials" (
    id int primary key generated always as identity
    , user_id uuid references users(id)
    , social social_network not null
    , link varchar(127) not null
    , is_preferred boolean not null
);

create table "work_places" (
    id int primary key generated always as identity
    , user_id uuid references users(id)
    , company_name varchar(31) not null
    , seniority varchar(31) not null
    , position varchar(31) not null
    , start_year smallint not null
    , end_year smallint
    , is_working_now boolean not null
);

create table "edu_places" (
    id int primary key generated always as identity
    , user_id uuid references users(id)
    , institution_name varchar(31) not null
    , grade edu_grade not null
    , level varchar(31)
    , specialization varchar(31) not null
    , start_year smallint not null
    , end_year smallint
    , is_studying_now boolean not null
);

-- +goose Down
drop table if exists "work_places";
drop table if exists "socials";
drop table if exists "users";
drop table if exists "soft_skills";
drop table if exists "key_skills";
drop table if exists "statuses";
drop type if exists "edu_grade";
drop type if exists "social_network";
drop type if exists "user_role";
