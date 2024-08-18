package repository

import (
	"context"
	"github.com/ukrainskykirill/chat-server/internal/model"
)

type ChatsRepository interface {
	Create(ctx context.Context, userIDs []int64) (int64, error)
	Delete(ctx context.Context, chatID int64) error
	IsExistsByUserIds(ctx context.Context, userIDs []int64) (bool, error)
	IsExistsById(ctx context.Context, chatID int64) (bool, error)
}

type MessagesRepository interface {
	Create(ctx context.Context, msgIn *model.MessageIn) error
}
