package model

type MessageIn struct {
	ChatID     int64
	UserID     int64
	ChatUserID int64
	Text       string
	Timestamp  string
}
