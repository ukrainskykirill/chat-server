package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/ukrainskykirill/chat-server/internal/model"
	"github.com/ukrainskykirill/chat-server/internal/repository"
	repoMocks "github.com/ukrainskykirill/chat-server/internal/repository/mocks"
	"github.com/ukrainskykirill/chat-server/internal/service/messages"
)

func TestSendMsg(t *testing.T) {
	t.Parallel()
	type chatsRepoMockFunc func(mc *minimock.Controller) repository.ChatsRepository
	type msgRepoMockFunc func(mc *minimock.Controller) repository.MessagesRepository

	type args struct {
		ctx   context.Context
		msgIn *model.MessageIn
	}

	var (
		ctx     = context.Background()
		mc      = minimock.NewController(t)
		repoErr = fmt.Errorf("repo error")

		id         = gofakeit.Int64()
		userID     = gofakeit.Int64()
		chatUserId = gofakeit.Int64()
		text       = gofakeit.BeerAlcohol()
		timestamp  = time.Now()

		msgIn = model.MessageIn{
			ChatID:     id,
			UserID:     userID,
			ChatUserID: chatUserId,
			Text:       text,
			Timestamp:  timestamp.String(),
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name          string
		args          args
		err           error
		chatsRepoMock chatsRepoMockFunc
		msgRepoMock   msgRepoMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:   ctx,
				msgIn: &msgIn,
			},
			err: nil,
			chatsRepoMock: func(mc *minimock.Controller) repository.ChatsRepository {
				mock := repoMocks.NewChatsRepositoryMock(mc)
				mock.GetChatUserIDMock.Expect(ctx, msgIn.ChatID, msgIn.UserID).Return(chatUserId, nil)
				return mock
			},
			msgRepoMock: func(mc *minimock.Controller) repository.MessagesRepository {
				mock := repoMocks.NewMessagesRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, &msgIn).Return(nil)
				return mock
			},
		},
		{
			name: "error case 1",
			args: args{
				ctx:   ctx,
				msgIn: &msgIn,
			},
			err: repoErr,
			chatsRepoMock: func(mc *minimock.Controller) repository.ChatsRepository {
				mock := repoMocks.NewChatsRepositoryMock(mc)
				mock.GetChatUserIDMock.Expect(ctx, msgIn.ChatID, msgIn.UserID).Return(0, repoErr)
				return mock
			},
			msgRepoMock: func(mc *minimock.Controller) repository.MessagesRepository {
				mock := repoMocks.NewMessagesRepositoryMock(mc)
				return mock
			},
		},
		{
			name: "error case 2",
			args: args{
				ctx:   ctx,
				msgIn: &msgIn,
			},
			err: repoErr,
			chatsRepoMock: func(mc *minimock.Controller) repository.ChatsRepository {
				mock := repoMocks.NewChatsRepositoryMock(mc)
				mock.GetChatUserIDMock.Expect(ctx, msgIn.ChatID, msgIn.UserID).Return(chatUserId, nil)
				return mock
			},
			msgRepoMock: func(mc *minimock.Controller) repository.MessagesRepository {
				mock := repoMocks.NewMessagesRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, &msgIn).Return(repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatsRepo := tt.chatsRepoMock(mc)
			msgRepo := tt.msgRepoMock(mc)

			service := messages.NewServ(
				msgRepo, chatsRepo,
			)

			err := service.Create(tt.args.ctx, tt.args.msgIn)
			require.Equal(t, tt.err, err)
		})
	}
}
