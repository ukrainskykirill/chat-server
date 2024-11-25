package service

import (
	"context"

	"github.com/ukrainskykirill/chat-server/internal/model"
)

type ChatsService interface {
	Create(ctx context.Context, userIDs []int64) (int64, error)
	Delete(ctx context.Context, chatID int64) error
	ConnectChat(connectInfo *model.ConnectInfo) error
}
type MessagesService interface {
	Create(ctx context.Context, msgIn *model.MessageIn) error
}

type StreamService interface {
	AddToChannel(ctx context.Context, chatID int64)
	AddMessageToChat(ctx context.Context, msgIn *model.MessageIn)
	GetChatChannel(ctx context.Context, chatID string) (chan *model.StreamMessage, error)
	AddStreamByUserID(connectInfo *model.ConnectInfo)
	SendMessageToStreams(chatID string, msg *model.StreamMessage) error
	DeleteStreamByUserID(chatID, userID string)
}
