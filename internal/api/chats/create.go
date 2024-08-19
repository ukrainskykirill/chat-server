package chats

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	prError "github.com/ukrainskykirill/chat-server/internal/error"
	gchat "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
)

func (i *Implementation) Create(ctx context.Context, req *gchat.CreateRequest) (*gchat.CreateResponse, error) {
	chatID, err := i.chatService.Create(ctx, req.UserIDs)
	if err != nil {
		switch {
		case errors.Is(err, prError.ErrChatAlreadyExist):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &gchat.CreateResponse{
		Id: chatID,
	}, nil
}
