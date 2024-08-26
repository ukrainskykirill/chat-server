package chats

import (
	"context"
)

func (s *chatServ) Create(ctx context.Context, userIDs []int64) (int64, error) {
	isRegular := len(userIDs) > 2

	var chatID int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		chatID, errTx = s.chatRepo.Create(ctx, isRegular)
		if errTx != nil {
			return errTx
		}

		errTx = s.chatRepo.CreateChatUsers(ctx, chatID, userIDs)
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return chatID, nil
}
