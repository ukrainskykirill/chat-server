package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"github.com/ukrainskykirill/platform_common/pkg/db"
	txMocks "github.com/ukrainskykirill/platform_common/pkg/db/mocks"

	"github.com/ukrainskykirill/chat-server/internal/repository"
	repoMocks "github.com/ukrainskykirill/chat-server/internal/repository/mocks"
	"github.com/ukrainskykirill/chat-server/internal/service/chats"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	type chatsRepoMockFunc func(mc *minimock.Controller) repository.ChatsRepository
	type msgRepoMockFunc func(mc *minimock.Controller) repository.MessagesRepository
	type txManagerMockFunc func(f db.Handler, mc *minimock.Controller) db.TxManager

	type args struct {
		ctx    context.Context
		chatId int64
	}

	var (
		ctx     = context.Background()
		mc      = minimock.NewController(t)
		repoErr = fmt.Errorf("repo error")

		chatId = gofakeit.Int64()
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name          string
		args          args
		err           error
		chatsRepoMock chatsRepoMockFunc
		msgRepoMock   msgRepoMockFunc
		txManagerMock txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:    ctx,
				chatId: chatId,
			},
			err: nil,
			chatsRepoMock: func(mc *minimock.Controller) repository.ChatsRepository {
				mock := repoMocks.NewChatsRepositoryMock(mc)
				mock.DeleteChatUsersMock.Expect(ctx, chatId).Return(nil)
				mock.DeleteChatMock.Expect(ctx, chatId).Return(nil)
				return mock
			},
			msgRepoMock: func(mc *minimock.Controller) repository.MessagesRepository {
				mock := repoMocks.NewMessagesRepositoryMock(mc)
				mock.DeleteByChatIDMock.Expect(ctx, chatId).Return(nil)
				return mock
			},
			txManagerMock: func(f db.Handler, mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "error case 1",
			args: args{
				ctx:    ctx,
				chatId: chatId,
			},
			err: repoErr,
			chatsRepoMock: func(mc *minimock.Controller) repository.ChatsRepository {
				mock := repoMocks.NewChatsRepositoryMock(mc)
				return mock
			},
			msgRepoMock: func(mc *minimock.Controller) repository.MessagesRepository {
				mock := repoMocks.NewMessagesRepositoryMock(mc)
				mock.DeleteByChatIDMock.Expect(ctx, chatId).Return(repoErr)
				return mock
			},
			txManagerMock: func(f db.Handler, mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "error case 2",
			args: args{
				ctx:    ctx,
				chatId: chatId,
			},
			err: repoErr,
			chatsRepoMock: func(mc *minimock.Controller) repository.ChatsRepository {
				mock := repoMocks.NewChatsRepositoryMock(mc)
				mock.DeleteChatUsersMock.Expect(ctx, chatId).Return(repoErr)
				return mock
			},
			msgRepoMock: func(mc *minimock.Controller) repository.MessagesRepository {
				mock := repoMocks.NewMessagesRepositoryMock(mc)
				mock.DeleteByChatIDMock.Expect(ctx, chatId).Return(nil)
				return mock
			},
			txManagerMock: func(f db.Handler, mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "error case 3",
			args: args{
				ctx:    ctx,
				chatId: chatId,
			},
			err: repoErr,
			chatsRepoMock: func(mc *minimock.Controller) repository.ChatsRepository {
				mock := repoMocks.NewChatsRepositoryMock(mc)
				mock.DeleteChatUsersMock.Expect(ctx, chatId).Return(nil)
				mock.DeleteChatMock.Expect(ctx, chatId).Return(repoErr)
				return mock
			},
			msgRepoMock: func(mc *minimock.Controller) repository.MessagesRepository {
				mock := repoMocks.NewMessagesRepositoryMock(mc)
				mock.DeleteByChatIDMock.Expect(ctx, chatId).Return(nil)
				return mock
			},
			txManagerMock: func(f db.Handler, mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})
				return mock
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatsRepo := tt.chatsRepoMock(mc)
			msgRepo := tt.msgRepoMock(mc)

			txManager := tt.txManagerMock(
				func(context.Context) error {
					err := tt.msgRepoMock(mc).DeleteByChatID(ctx, chatId)
					if err != nil {
						return err
					}
					err = tt.chatsRepoMock(mc).DeleteChatUsers(ctx, chatId)
					if err != nil {
						return err
					}
					err = tt.chatsRepoMock(mc).DeleteChat(ctx, chatId)
					if err != nil {
						return err
					}
					return nil
				}, mc,
			)

			service := chats.NewServ(
				txManager, chatsRepo, msgRepo,
			)

			err := service.Delete(tt.args.ctx, tt.args.chatId)
			require.Equal(t, tt.err, err)
		})
	}
}
