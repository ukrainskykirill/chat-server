-- +goose Up
-- +goose StatementBegin
CREATE TABLE chats (
    id serial PRIMARY KEY,
    is_regular boolean
);
CREATE TABLE chat_users (
    id serial PRIMARY KEY,
    chat_id integer,
    username varchar(500)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chats;
DROP TABLE chat_users;
-- +goose StatementEnd