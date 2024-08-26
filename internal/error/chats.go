package error

import "errors"

var (
	ErrChatNotFound         = errors.New("chat not found")
	ErrChatUserNotFound     = errors.New("chat user not found")
	ErrUserIsNotParticipant = errors.New("user is not participant")
	ErrChatAlreadyExist     = errors.New("chat already exist")
)
