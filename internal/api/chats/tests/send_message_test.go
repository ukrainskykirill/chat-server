package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukrainskykirill/chat-server/internal/api/chats"
	"github.com/ukrainskykirill/chat-server/internal/model"
	"github.com/ukrainskykirill/chat-server/internal/service"
	serviceMocks "github.com/ukrainskykirill/chat-server/internal/service/mocks"
	desc "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()
	type chatsServiceMockFunc func(mc *minimock.Controller) service.ChatsService
	type msgServiceMockFunc func(mc *minimock.Controller) service.MessagesService

	type args struct {
		ctx context.Context
		req *desc.SendMessageRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		chatId    = gofakeit.Int64()
		userId    = gofakeit.Int64()
		text      = gofakeit.BeerName()
		timestamp = timestamppb.Now()

		serviceErr = fmt.Errorf("service error")

		req = &desc.SendMessageRequest{
			ChatID:    chatId,
			UserId:    userId,
			Text:      text,
			Timestamp: timestamp,
		}

		messageIn = &model.MessageIn{
			ChatID:    chatId,
			UserID:    userId,
			Text:      text,
			Timestamp: timestamp.AsTime().Format(time.UnixDate),
		}

		res = &emptypb.Empty{}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name             string
		args             args
		want             *emptypb.Empty
		err              error
		chatsServiceMock chatsServiceMockFunc
		msgServiceMock   msgServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			chatsServiceMock: func(mc *minimock.Controller) service.ChatsService {
				mock := serviceMocks.NewChatsServiceMock(mc)
				return mock
			},
			msgServiceMock: func(mc *minimock.Controller) service.MessagesService {
				mock := serviceMocks.NewMessagesServiceMock(mc)
				mock.CreateMock.Expect(ctx, messageIn).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  status.Error(codes.Internal, serviceErr.Error()),
			chatsServiceMock: func(mc *minimock.Controller) service.ChatsService {
				mock := serviceMocks.NewChatsServiceMock(mc)
				return mock
			},
			msgServiceMock: func(mc *minimock.Controller) service.MessagesService {
				mock := serviceMocks.NewMessagesServiceMock(mc)
				mock.CreateMock.Expect(ctx, messageIn).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatsServiceMock := tt.chatsServiceMock(mc)
			msgServiceMock := tt.msgServiceMock(mc)

			api := chats.NewChatImplementation(
				chatsServiceMock, msgServiceMock,
			)

			chatID, err := api.SendMessage(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, chatID)
		})
	}
}
