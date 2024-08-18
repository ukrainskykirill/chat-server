package chats

import (
	"context"

	gchat "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Delete(ctx context.Context, req *gchat.DeleteRequest) (*emptypb.Empty, error) {
	err := i.chatService.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
