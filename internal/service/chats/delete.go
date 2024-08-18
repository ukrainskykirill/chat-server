package chats

import "context"

func (s *chatServ) Delete(ctx context.Context, chatID int64) error {
	err := s.repo.Delete(ctx, chatID)
	if err != nil {
		return err
	}
	return err
}
