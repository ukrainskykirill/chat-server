package chats

import (
	"context"
	// "fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gchat "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
)

func (i *Implementation) Create(ctx context.Context, req *gchat.CreateRequest) (*gchat.CreateResponse, error) {
	// if req.UserIDs[0] == 1 {
	// 	return nil, status.Error(codes.Internal, fmt.Errorf("invalid req %d", req.UserIDs[0]).Error())
	// }
	chatID, err := i.chatService.Create(ctx, req.UserIDs)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gchat.CreateResponse{
		Id: chatID,
	}, nil
}
