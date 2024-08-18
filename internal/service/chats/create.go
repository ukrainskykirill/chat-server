package chats

import "context"

func (s *chatServ) Create(ctx context.Context, userIDs []int64) (int64, error) {
	chatID, err := s.repo.Create(ctx, userIDs)
	if err != nil {
		return 0, err
	}
	return chatID, nil
}
