package chats

import (
	"github.com/ukrainskykirill/platform_common/pkg/db"

	"github.com/ukrainskykirill/chat-server/internal/repository"
	"github.com/ukrainskykirill/chat-server/internal/service"
)

type chatServ struct {
	txManager     db.TxManager
	chatRepo      repository.ChatsRepository
	msgRepo       repository.MessagesRepository
	streamService service.StreamService
}

func NewServ(
	txManager db.TxManager,
	chatsRepo repository.ChatsRepository,
	msgRepo repository.MessagesRepository,
	streamService service.StreamService,
) service.ChatsService {
	return &chatServ{
		txManager:     txManager,
		chatRepo:      chatsRepo,
		msgRepo:       msgRepo,
		streamService: streamService,
	}
}
