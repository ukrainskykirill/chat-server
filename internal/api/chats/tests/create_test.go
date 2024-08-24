package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukrainskykirill/chat-server/internal/api/chats"
	"github.com/ukrainskykirill/chat-server/internal/service"
	serviceMocks "github.com/ukrainskykirill/chat-server/internal/service/mocks"
	desc "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type chatsServiceMockFunc func(mc *minimock.Controller) service.ChatsService
	type msgServiceMockFunc func(mc *minimock.Controller) service.MessagesService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id      = gofakeit.Int64()
		userIDs = []int64{gofakeit.Int64(), gofakeit.Int64()}

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateRequest{
			UserIDs: userIDs,
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name             string
		args             args
		want             *desc.CreateResponse
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
				mock.CreateMock.Expect(ctx, req.UserIDs).Return(id, nil)
				return mock
			},
			msgServiceMock: func(mc *minimock.Controller) service.MessagesService {
				mock := serviceMocks.NewMessagesServiceMock(mc)
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
				mock.CreateMock.Expect(ctx, req.UserIDs).Return(0, serviceErr)
				return mock
			},
			msgServiceMock: func(mc *minimock.Controller) service.MessagesService {
				mock := serviceMocks.NewMessagesServiceMock(mc)
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

			chatID, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, chatID)
		})
	}
}
