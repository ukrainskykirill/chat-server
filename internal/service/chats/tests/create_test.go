package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/ukrainskykirill/chat-server/internal/client/db"
	txMocks "github.com/ukrainskykirill/chat-server/internal/client/db/mocks"
	"github.com/ukrainskykirill/chat-server/internal/repository"
	repoMocks "github.com/ukrainskykirill/chat-server/internal/repository/mocks"
	"github.com/ukrainskykirill/chat-server/internal/service/chats"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type chatsRepoMockFunc func(mc *minimock.Controller) repository.ChatsRepository
	type msgRepoMockFunc func(mc *minimock.Controller) repository.MessagesRepository
	type txManagerMockFunc func(f db.Handler, mc *minimock.Controller) db.TxManager

	type args struct {
		ctx       context.Context
		userIDs   []int64
		isRegular bool
	}

	var (
		ctx     = context.Background()
		mc      = minimock.NewController(t)
		repoErr = fmt.Errorf("repo error")

		id        = gofakeit.Int64()
		userIDs   = []int64{gofakeit.Int64(), gofakeit.Int64()}
		isRegular = true
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name          string
		args          args
		want          int64
		err           error
		chatsRepoMock chatsRepoMockFunc
		msgRepoMock   msgRepoMockFunc
		txManagerMock txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:       ctx,
				userIDs:   userIDs,
				isRegular: isRegular,
			},
			want: id,
			err:  nil,
			chatsRepoMock: func(mc *minimock.Controller) repository.ChatsRepository {
				mock := repoMocks.NewChatsRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, isRegular).Return(id, nil)
				mock.CreateChatUsersMock.Expect(ctx, id, userIDs).Return(nil)
				return mock
			},
			msgRepoMock: func(mc *minimock.Controller) repository.MessagesRepository {
				mock := repoMocks.NewMessagesRepositoryMock(mc)
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
				ctx:       ctx,
				userIDs:   userIDs,
				isRegular: isRegular,
			},
			want: 0,
			err:  repoErr,
			chatsRepoMock: func(mc *minimock.Controller) repository.ChatsRepository {
				mock := repoMocks.NewChatsRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, isRegular).Return(0, repoErr)
				return mock
			},
			msgRepoMock: func(mc *minimock.Controller) repository.MessagesRepository {
				mock := repoMocks.NewMessagesRepositoryMock(mc)
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
				ctx:       ctx,
				userIDs:   userIDs,
				isRegular: isRegular,
			},
			want: 0,
			err:  repoErr,
			chatsRepoMock: func(mc *minimock.Controller) repository.ChatsRepository {
				mock := repoMocks.NewChatsRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, isRegular).Return(id, nil)
				mock.CreateChatUsersMock.Expect(ctx, id, userIDs).Return(repoErr)
				return mock
			},
			msgRepoMock: func(mc *minimock.Controller) repository.MessagesRepository {
				mock := repoMocks.NewMessagesRepositoryMock(mc)
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
					chatID, err := tt.chatsRepoMock(mc).Create(ctx, isRegular)
					if err != nil {
						return err
					}
					err = tt.chatsRepoMock(mc).CreateChatUsers(ctx, chatID, userIDs)
					if err != nil {
						return err
					}
					return nil
				}, mc,
			)

			service := chats.NewServ(
				txManager, chatsRepo, msgRepo,
			)

			id, err := service.Create(tt.args.ctx, tt.args.userIDs)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, id)
		})
	}
}
