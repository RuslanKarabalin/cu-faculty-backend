-- +goose Up
insert into "statuses" (content) values ('Собираю команду');
insert into "statuses" (content) values ('Хочу нетворкинг');
insert into "statuses" (content) values ('Ищу стартап');
insert into "statuses" (content) values ('Работаю');
insert into "statuses" (content) values ('Устал');

-- +goose Down
truncate table "statuses";
