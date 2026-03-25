-- +goose Up
-- insert into "statuses"(content) values('');
-- insert into "key_skills"(name) values('');
-- insert into "soft_skills"(name) values('');
-- insert into "universities"(name) values('');
-- insert into "faqs"(question, answer) values('', '');

-- +goose Down
truncate table "faqs";
truncate table "universities";
truncate table "soft_skills";
truncate table "key_skills";
truncate table "statuses";
