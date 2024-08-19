package chats

import (
	"github.com/ukrainskykirill/chat-server/internal/client/db"
	"github.com/ukrainskykirill/chat-server/internal/repository"
	"github.com/ukrainskykirill/chat-server/internal/service"
)

type chatServ struct {
	txManager db.TxManager
	chatRepo  repository.ChatsRepository
	msgRepo   repository.MessagesRepository
}

func NewServ(txManager db.TxManager, chatsRepo repository.ChatsRepository, msgRepo repository.MessagesRepository) service.ChatsService {
	return &chatServ{
		txManager: txManager,
		chatRepo:  chatsRepo,
		msgRepo:   msgRepo,
	}
}
