package converter

import (
	"time"

	"github.com/ukrainskykirill/chat-server/internal/model"
	gchat "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
)

func ToModelMessageIn(msgIn *gchat.SendMessageRequest) *model.MessageIn {
	timestamp := msgIn.Timestamp.AsTime().Format(time.UnixDate)

	return &model.MessageIn{
		ChatID:    1,
		Text:      msgIn.Text,
		Timestamp: timestamp,
	}
}
