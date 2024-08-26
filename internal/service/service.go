package service

import (
	"context"

	"github.com/ukrainskykirill/chat-server/internal/model"
)

type ChatsService interface {
	Create(ctx context.Context, userIDs []int64) (int64, error)
	Delete(ctx context.Context, chatID int64) error
}
type MessagesService interface {
	Create(ctx context.Context, msgIn *model.MessageIn) error
}
