package chats

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	prError "github.com/ukrainskykirill/chat-server/internal/error"
	gchat "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
)

func (i *Implementation) Delete(ctx context.Context, req *gchat.DeleteRequest) (*emptypb.Empty, error) {
	err := i.chatService.Delete(ctx, req.Id)
	if err != nil {
		switch {
		case errors.Is(err, prError.ErrChatNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &emptypb.Empty{}, nil
}
