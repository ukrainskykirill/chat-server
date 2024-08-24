package messages

import (
	"context"

	"github.com/ukrainskykirill/chat-server/internal/model"
)

func (s *messagesServ) Create(ctx context.Context, msgIn *model.MessageIn) error {
	chatUserID, err := s.chatsRepo.GetChatUserID(ctx, msgIn.ChatID, msgIn.UserID)
	if err != nil {
		return err
	}
	msgIn.ChatUserID = chatUserID

	err = s.msgRepo.Create(ctx, msgIn)
	if err != nil {
		return err
	}
	return nil
}
