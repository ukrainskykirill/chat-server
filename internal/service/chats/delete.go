package chats

import (
	"context"
	"fmt"

	prError "github.com/ukrainskykirill/chat-server/internal/error"
)

func (s *chatServ) Delete(ctx context.Context, chatID int64) error {
	isExist, err := s.chatRepo.IsExistsByID(ctx, chatID)
	if err != nil {
		return err
	}
	if !isExist {
		return fmt.Errorf("chats.Delete - %w", prError.ErrChatNotFound)
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		errTx = s.msgRepo.DeleteByChatID(ctx, chatID)
		if errTx != nil {
			return errTx
		}

		errTx = s.chatRepo.DeleteChatUsers(ctx, chatID)
		if errTx != nil {
			return errTx
		}

		errTx = s.chatRepo.DeleteChat(ctx, chatID)
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
