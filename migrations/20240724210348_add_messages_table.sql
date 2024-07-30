-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
    id serial PRIMARY KEY,
    chat_user_id integer,
    text varchar(5000),
    timestamp timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE messages;
-- +goose StatementEnd
