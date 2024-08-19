package repository

import (
	"context"

	"github.com/ukrainskykirill/chat-server/internal/model"
)

type ChatsRepository interface {
	Create(ctx context.Context, isRegular bool) (int64, error)
	GetChatUserID(ctx context.Context, chatID int64, userID int64) (int64, error)
	CreateChatUsers(ctx context.Context, chatID int64, userIDs []int64) error
	DeleteChat(ctx context.Context, chatID int64) error
	DeleteChatUsers(ctx context.Context, chatID int64) error
	IsExistsByID(ctx context.Context, chatID int64) (bool, error)
	IsUserParticipant(ctx context.Context, chatID int64, userID int64) (bool, error)
}

type MessagesRepository interface {
	Create(ctx context.Context, msgIn *model.MessageIn) error
	DeleteByChatID(ctx context.Context, chatID int64) error
}
