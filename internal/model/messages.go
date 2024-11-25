package model

import (
	gchat "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
)

type MessageIn struct {
	ChatID     int64
	UserID     int64
	ChatUserID int64
	Text       string
	Timestamp  string
}

type StreamMessage struct {
	UserID    int64
	Text      string
	Timestamp string
}

type ConnectInfo struct {
	UserID string
	ChatID string
	Stream gchat.ChatV1_ConnectChatServer
}
