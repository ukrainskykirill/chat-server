package messages

import (
	"github.com/ukrainskykirill/chat-server/internal/repository"
	"github.com/ukrainskykirill/chat-server/internal/service"
)

type messagesServ struct {
	chatsRepo repository.ChatsRepository
	msgRepo   repository.MessagesRepository
}

func NewServ(msgRepo repository.MessagesRepository, chatsRepo repository.ChatsRepository) service.MessagesService {
	return &messagesServ{
		msgRepo:   msgRepo,
		chatsRepo: chatsRepo,
	}
}
