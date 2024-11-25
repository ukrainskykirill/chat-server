package messages

import (
	"github.com/ukrainskykirill/chat-server/internal/repository"
	"github.com/ukrainskykirill/chat-server/internal/service"
)

type messagesServ struct {
	chatsRepo     repository.ChatsRepository
	msgRepo       repository.MessagesRepository
	streamService service.StreamService
}

func NewServ(
	msgRepo repository.MessagesRepository,
	chatsRepo repository.ChatsRepository,
	streamService service.StreamService,
) service.MessagesService {
	return &messagesServ{
		msgRepo:       msgRepo,
		chatsRepo:     chatsRepo,
		streamService: streamService,
	}
}
