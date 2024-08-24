package chats

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ukrainskykirill/chat-server/internal/converter"
	prError "github.com/ukrainskykirill/chat-server/internal/error"
	gchat "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
)

func (i *Implementation) SendMessage(ctx context.Context, req *gchat.SendMessageRequest) (*emptypb.Empty, error) {
	err := i.msgService.Create(ctx, converter.ToModelMessageIn(req))
	if err != nil {
		switch {
		case errors.Is(err, prError.ErrChatNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, prError.ErrChatUserNotFound):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &emptypb.Empty{}, nil
}
