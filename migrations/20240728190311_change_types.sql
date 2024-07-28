-- +goose Up
-- +goose StatementBegin
ALTER TABLE messages ALTER COLUMN text TYPE TEXT;
ALTER TABLE messages ALTER COLUMN text SET NOT NULL;
ALTER TABLE chat_users ALTER COLUMN username TYPE TEXT;
ALTER TABLE chat_users ALTER COLUMN username SET NOT NULL;
ALTER TABLE messages ALTER COLUMN timestamp SET NOT NULL;
ALTER TABLE messages ALTER COLUMN chat_user_id SET NOT NULL;
ALTER TABLE chat_users ALTER COLUMN chat_id SET NOT NULL;
ALTER TABLE chats ALTER COLUMN is_regular SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE messages ALTER COLUMN text DROP NOT NULL;
ALTER TABLE chat_users ALTER COLUMN username DROP NOT NULL;
ALTER TABLE messages ALTER COLUMN text TYPE varchar(5000);
ALTER TABLE chat_users ALTER COLUMN username TYPE varchar(500);
ALTER TABLE chat_users ALTER COLUMN chat_id DROP NOT NULL;
ALTER TABLE chats ALTER COLUMN is_regular DROP NOT NULL;
ALTER TABLE messages ALTER COLUMN timestamp DROP NOT NULL;
ALTER TABLE messages ALTER COLUMN chat_user_id DROP NOT NULL;
-- +goose StatementEnd
