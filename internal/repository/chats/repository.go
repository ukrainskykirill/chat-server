package chats

import (
	"context"
	"fmt"

	"github.com/ukrainskykirill/chat-server/internal/client/db"
	"github.com/ukrainskykirill/chat-server/internal/repository"
)

const (
	chatsRepo = "chats_repository"
)

type repo struct {
	db db.Client
}

func NewChatRepository(db db.Client) repository.ChatsRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, userIDs []int64) (int64, error) {
	fmt.Println("create repo")
	return 0, nil
}

func (r *repo) Delete(ctx context.Context, chatID int64) error {
	return nil
}

func (r *repo) IsExistsByUserIds(ctx context.Context, userIDs []int64) (bool, error) {
	return true, nil
}

func (r *repo) IsExistsById(ctx context.Context, chatID int64) (bool, error) {
	return true, nil
}
