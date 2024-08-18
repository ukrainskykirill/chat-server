package chats

import (
	"context"

	"github.com/ukrainskykirill/chat-server/internal/converter"
	gchat "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) SendMessage(ctx context.Context, req *gchat.SendMessageRequest) (*emptypb.Empty, error) {
	err := i.msgService.Create(ctx, converter.ToModelMessageIn(req))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
