package converter

import (
	"fmt"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukrainskykirill/chat-server/internal/model"
	gchat "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
)

func ToModelConnectChat(req *gchat.ConnectChatRequest, stream gchat.ChatV1_ConnectChatServer) *model.ConnectInfo {
	return &model.ConnectInfo{
		ChatID: req.GetChatId(),
		UserID: req.GetUserId(),
		Stream: stream,
	}
}

func ToGMessageFromModel(msg *model.StreamMessage) *gchat.Message {
	timestamp, err := time.Parse(time.RFC3339, msg.Timestamp)
	if err != nil {
		fmt.Println("cant cast string to time")
	}

	return &gchat.Message{
		UserId:    strconv.Itoa(int(msg.UserID)),
		Text:      msg.Text,
		Timestamp: timestamppb.New(timestamp),
	}
}
