package chats

import (
	"context"

	"github.com/ukrainskykirill/chat-server/internal/client/db"
	"github.com/ukrainskykirill/chat-server/internal/model"
	"github.com/ukrainskykirill/chat-server/internal/repository"
)

const (
	chatsRepo = "messages_repository"
)

type repo struct {
	db db.Client
}

func NewMessageRepository(db db.Client) repository.MessagesRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, msgIn *model.MessageIn) error {
	return nil
}
