-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
    id serial PRIMARY KEY,
    chat_user_id integer NOT NULL,
    text TEXT NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE messages;
-- +goose StatementEnd
