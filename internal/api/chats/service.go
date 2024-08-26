package chats

import (
	"github.com/ukrainskykirill/chat-server/internal/service"
	gchat "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
)

type Implementation struct {
	gchat.UnimplementedChatV1Server
	chatService service.ChatsService
	msgService  service.MessagesService
}

func NewChatImplementation(chat service.ChatsService, msg service.MessagesService) *Implementation {
	return &Implementation{
		chatService: chat,
		msgService:  msg,
	}
}
