package chats

import (
	"context"

	gchat "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Create(ctx context.Context, req *gchat.CreateRequest) (*gchat.CreateResponse, error) {
	chatID, err := i.chatService.Create(ctx, []int64{1, 2})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gchat.CreateResponse{
		Id: chatID,
	}, nil
}
