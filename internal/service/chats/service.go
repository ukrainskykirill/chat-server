package chats

import (
	"github.com/ukrainskykirill/chat-server/internal/repository"
	"github.com/ukrainskykirill/chat-server/internal/service"
)

type chatServ struct {
	repo repository.ChatsRepository
}

func NewServ(chatsRepo repository.ChatsRepository) service.ChatsService {
	return &chatServ{
		repo: chatsRepo,
	}
}
