package chats

import (
	"context"
	"time"

	"github.com/ukrainskykirill/platform_common/pkg/db"

	"github.com/ukrainskykirill/chat-server/internal/model"
	"github.com/ukrainskykirill/chat-server/internal/repository"
)

const (
	msgRepo   = "messages_repository"
	createMsg = msgRepo + "." + "CreateMsg"
	deleteMsg = msgRepo + "." + "DeleteChat"
)

type repo struct {
	db db.Client
}

func NewMessageRepository(db db.Client) repository.MessagesRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, msgIn *model.MessageIn) error {
	rowSQL := `INSERT INTO messages (text, chat_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4)`

	q := db.Query{
		Name:     createMsg,
		QueryRaw: rowSQL,
	}

	_, err := r.db.DB().ExecContext(
		ctx,
		q,
		msgIn.Text, msgIn.ChatUserID, time.Now(), time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) DeleteByChatID(ctx context.Context, chatID int64) error {
	rowSQL := `DELETE FROM messages WHERE chat_user_id IN (SELECT id FROM chat_users WHERE chat_id = $1);`

	q := db.Query{
		Name:     deleteMsg,
		QueryRaw: rowSQL,
	}

	_, err := r.db.DB().ExecContext(
		ctx,
		q,
		chatID,
	)
	if err != nil {
		return err
	}
	return nil
}
