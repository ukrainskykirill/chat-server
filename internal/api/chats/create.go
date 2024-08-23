package chats

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gchat "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
)

func (i *Implementation) Create(ctx context.Context, req *gchat.CreateRequest) (*gchat.CreateResponse, error) {
	chatID, err := i.chatService.Create(ctx, req.UserIDs)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gchat.CreateResponse{
		Id: chatID,
	}, nil
}
