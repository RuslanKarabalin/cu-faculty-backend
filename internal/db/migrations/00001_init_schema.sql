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

create type "work_grade" as enum (
    'intern'
    , 'junior'
    , 'junior_plus'
    , 'middle'
    , 'middle_plus'
    , 'senior'
    , 'staff'
    , 'principal'
    , 'lead'
    , 'head'
);

create type "event_category" as enum (
    'networking'
    , 'professional'
    , 'partner'
);

create table "statuses" (
    id int primary key generated always as identity
    , content varchar(127) not null
);

create table "key_skills" (
    id int primary key generated always as identity
    , name varchar(63) not null
);

create table "soft_skills" (
    id int primary key generated always as identity
    , name varchar(63) not null
);

create table "companies" (
    id int primary key generated always as identity
    , name varchar(63) not null
);

create table "work_positions" (
    id int primary key generated always as identity
    , name varchar(63) not null
);

create table "universities" (
    id int primary key generated always as identity
    , name varchar(127) not null
    , short_name varchar(31) not null
);

create table "faqs" (
    id int primary key generated always as identity
    , question varchar(255) not null
    , answer varchar(255) not null
);

create table "users" (
    id uuid primary key default gen_random_uuid()
    , email varchar(255) not null
    , photo_s3_key varchar(255)
    , first_name varchar(31) not null
    , last_name varchar(31) not null
    , bio varchar(255)
    , birth_date date not null
    , speciality varchar(63)
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
    , company_name varchar(63) not null
    , grade work_grade not null
    , position varchar(63) not null
    , start_year smallint not null
    , end_year smallint
    , is_working_now boolean not null
);

create table "edu_places" (
    id int primary key generated always as identity
    , user_id uuid references users(id)
    , university_id int not null references universities(id)
    , grade edu_grade not null
    , level varchar(31)
    , specialization varchar(63) not null
    , start_year smallint not null
    , end_year smallint
    , is_studying_now boolean not null
);

create table "announcements" (
    id uuid primary key default gen_random_uuid()
    , author_id uuid references users(id)
    , title varchar(63) not null
    , content varchar(255) not null
    , created_at timestamptz not null default now()
    , is_archived boolean not null
);

create table "news" (
    id uuid primary key default gen_random_uuid()
    , author_id uuid references users(id)
    , photo_s3_key varchar(255) not null
    , title varchar(31) not null
    , content varchar(255) not null
    , publish_days smallint not null default 7
    , is_draft boolean not null default false
    , created_at timestamptz not null default now()
);

create table "events" (
    id uuid primary key default gen_random_uuid()
    , author_id uuid references users(id)
    , photo_s3_key varchar(255) not null
    , title varchar(31) not null
    , content varchar(255) not null
    , place varchar(63) not null
    , category event_category not null
    , starts_at timestamptz not null
    , registration_link varchar(255)
    , is_draft boolean not null default false
    , created_at timestamptz not null default now()
);

create table "user_key_skills" (
    user_id uuid references users(id)
    , key_skill_id int references key_skills(id)
    , primary key (user_id, key_skill_id)
);

create table "user_soft_skills" (
    user_id uuid references users(id)
    , soft_skill_id int references soft_skills(id)
    , primary key (user_id, soft_skill_id)
);

create table "contacts" (
    user_id uuid references users(id)
    , contact_id uuid references users(id)
    , note varchar(255)
    , primary key (user_id, contact_id)
);

create table "saved_users" (
    user_id uuid references users(id)
    , saved_user_id uuid references users(id)
    , primary key (user_id, saved_user_id)
);

create table "announcement_responses" (
    user_id uuid references users(id)
    , announcement_id uuid references announcements(id)
    , primary key (user_id, announcement_id)
);

create table "event_responses" (
    user_id uuid references users(id)
    , event_id uuid references events(id)
    , primary key (user_id, event_id)
);

-- +goose Down
drop table if exists "event_responses";
drop table if exists "announcement_responses";
drop table if exists "saved_users";
drop table if exists "contacts";
drop table if exists "user_soft_skills";
drop table if exists "user_key_skills";
drop table if exists "events";
drop table if exists "news";
drop table if exists "announcements";
drop table if exists "edu_places";
drop table if exists "work_places";
drop table if exists "socials";
drop table if exists "users";
drop table if exists "faqs";
drop table if exists "universities";
drop table if exists "work_positions";
drop table if exists "companies";
drop table if exists "soft_skills";
drop table if exists "key_skills";
drop table if exists "statuses";

drop type if exists "event_category";
drop type if exists "work_grade";
drop type if exists "edu_grade";
drop type if exists "social_network";
drop type if exists "user_role";
