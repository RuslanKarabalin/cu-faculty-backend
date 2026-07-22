-- +goose Up
insert into "statuses" (content) values ('Ищу проект');
insert into "statuses" (content) values ('Собираю команду');
insert into "statuses" (content) values ('Хочу нетворкинг');
insert into "statuses" (content) values ('Работаю');
insert into "statuses" (content) values ('Хочу на отдых');
insert into "statuses" (content) values ('Не беспокоить');

-- +goose Down
truncate table "statuses";
