package chats

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukrainskykirill/chat-server/internal/converter"
	gchat "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
)

func (i *Implementation) ConnectChat(req *gchat.ConnectChatRequest, stream gchat.ChatV1_ConnectChatServer) error {

	err := i.chatService.ConnectChat(converter.ToModelConnectChat(req, stream))
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}
