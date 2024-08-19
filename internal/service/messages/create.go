package messages

import (
	"context"
	"fmt"

	prError "github.com/ukrainskykirill/chat-server/internal/error"
	"github.com/ukrainskykirill/chat-server/internal/model"
)

func (s *messagesServ) Create(ctx context.Context, msgIn *model.MessageIn) error {
	isExist, err := s.chatsRepo.IsExistsByID(ctx, msgIn.ChatID)
	if err != nil {
		return err
	}
	if !isExist {
		return fmt.Errorf("msg service.Create - %w", prError.ErrChatNotFound)
	}

	isParticipant, err := s.chatsRepo.IsUserParticipant(ctx, msgIn.ChatID, msgIn.UserID)
	if err != nil {
		return err
	}
	if !isParticipant {
		return fmt.Errorf("msg service.Create - %w", prError.ErrUserIsNotParticipant)
	}

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
