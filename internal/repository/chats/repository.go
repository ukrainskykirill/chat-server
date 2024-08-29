package chats

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/ukrainskykirill/platform_common/pkg/db"

	prError "github.com/ukrainskykirill/chat-server/internal/error"
	"github.com/ukrainskykirill/chat-server/internal/repository"
)

const (
	chatsRepo                = "chats_repository"
	chatUsersTable           = "chat_users"
	chatUsersUserIDColumn    = "user_id"
	chatUsersChatIDColumn    = "chat_id"
	chatUsersCreatedAtColumn = "created_at"
	chatUsersUpdatedAtColumn = "updated_at"
	createChat               = chatsRepo + "." + "CreateChat"
	deleteChat               = chatsRepo + "." + "DeleteChat"
	deleteChatUsers          = chatsRepo + "." + "DeleteChatUsers"
	getChatUserID            = chatsRepo + "." + "GetChatUserId"
	isUserParticipant        = chatsRepo + "." + "IsUserParticipant"
)

type repo struct {
	db db.Client
}

func NewChatRepository(db db.Client) repository.ChatsRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, isRegular bool) (int64, error) {
	var chatID int64

	rowSQL := `INSERT INTO chats (is_regular, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id;`

	q := db.Query{
		Name:     createChat,
		QueryRaw: rowSQL,
	}

	err := r.db.DB().QueryRowContext(
		ctx,
		q,
		isRegular, time.Now(), time.Now(),
	).Scan(&chatID)
	if err != nil {
		return 0, err
	}

	fmt.Println(color.BlueString("Create chat: is regular - %b, with ctx: %v", isRegular, ctx))
	return chatID, nil
}

func (r *repo) GetChatUserID(ctx context.Context, chatID int64, userID int64) (int64, error) {
	var chatUserID int64

	q := db.Query{
		Name:     getChatUserID,
		QueryRaw: `SELECT id FROM chat_users WHERE chat_id = $1 and user_id = $2;`,
	}

	err := r.db.DB().QueryRowContext(
		ctx,
		q,
		chatID,
		userID,
	).Scan(&chatUserID)
	if err != nil {
		if errors.Is(err, prError.ErrChatUserNotFound) {
			return 0, prError.ErrChatUserNotFound
		}
		return 0, err
	}

	return chatUserID, nil

}

func (r *repo) CreateChatUsers(ctx context.Context, chatID int64, userIDs []int64) error {
	rows := [][]any{}
	for _, username := range userIDs {
		newRow := []any{username, chatID, time.Now(), time.Now()}
		rows = append(rows, newRow)
	}
	_, err := r.db.DB().CopyFrom(
		ctx,
		chatUsersTable,
		[]string{
			chatUsersUserIDColumn,
			chatUsersChatIDColumn,
			chatUsersCreatedAtColumn,
			chatUsersUpdatedAtColumn,
		},
		rows,
	)

	if err != nil {
		return err
	}
	return nil
}

func (r *repo) DeleteChat(ctx context.Context, chatID int64) error {
	rowSQL := `DELETE FROM chats WHERE id = $1`

	q := db.Query{
		Name:     deleteChatUsers,
		QueryRaw: rowSQL,
	}

	tag, err := r.db.DB().ExecContext(
		ctx,
		q,
		chatID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return prError.ErrChatNotFound
	}

	return nil
}

func (r *repo) DeleteChatUsers(ctx context.Context, chatID int64) error {
	rowSQL := `DELETE FROM chat_users WHERE chat_id = $1`

	q := db.Query{
		Name:     deleteChatUsers,
		QueryRaw: rowSQL,
	}

	_, err := r.db.DB().ExecContext(
		ctx,
		q,
		chatID,
	)
	if err != nil {
		return err
	}

	return nil
}
