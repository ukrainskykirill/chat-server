-- +goose Up
-- +goose StatementBegin
CREATE TABLE chats (
    id serial PRIMARY KEY,
    is_regular boolean NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);
CREATE TABLE chat_users (
    id serial PRIMARY KEY,
    chat_id integer NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chats;
DROP TABLE chat_users;
-- +goose StatementEnd