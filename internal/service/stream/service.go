package stream

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/ukrainskykirill/chat-server/internal/converter"
	"github.com/ukrainskykirill/chat-server/internal/model"
	"github.com/ukrainskykirill/chat-server/internal/service"
	desc "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
)

type Chat struct {
	streams  map[string]desc.ChatV1_ConnectChatServer
	mxStream sync.RWMutex
}

type streamService struct {
	сhats  map[string]*Chat
	mxChat sync.RWMutex

	channels  map[string]chan *model.StreamMessage
	mxChannel sync.RWMutex
}

func NewStreamService() service.StreamService {
	return &streamService{
		сhats:    make(map[string]*Chat),
		channels: make(map[string]chan *model.StreamMessage),
	}
}

func (s *streamService) AddToChannel(ctx context.Context, chatID int64) {
	s.channels[strconv.FormatInt(chatID, 10)] = make(chan *model.StreamMessage, 100)
}

func (s *streamService) AddMessageToChat(ctx context.Context, msgIn *model.MessageIn) {
	s.mxChannel.RLock()
	chatChan, ok := s.channels[strconv.FormatInt(msgIn.ChatID, 10)]
	s.mxChannel.RUnlock()

	if !ok {
		fmt.Println("chat not found")
	}

	chatChan <- &model.StreamMessage{
		UserID:    msgIn.UserID,
		Text:      msgIn.Text,
		Timestamp: msgIn.Text,
	}

}

func (s *streamService) GetChatChannel(ctx context.Context, chatID string) (chan *model.StreamMessage, error) {
	s.mxChannel.RLock()
	chatChan, ok := s.channels[chatID]
	s.mxChannel.RUnlock()

	if !ok {
		return nil, fmt.Errorf("chat didnt find by this id: %s", chatID)
	}

	return chatChan, nil
}

func (s *streamService) AddStreamByUserID(connectInfo *model.ConnectInfo) {
	s.mxChat.Lock()
	if _, okChat := s.сhats[connectInfo.ChatID]; !okChat {
		s.сhats[connectInfo.ChatID] = &Chat{
			streams: make(map[string]desc.ChatV1_ConnectChatServer),
		}
	}
	s.mxChat.Unlock()

	s.сhats[connectInfo.ChatID].mxStream.Lock()
	s.сhats[connectInfo.ChatID].streams[connectInfo.UserID] = connectInfo.Stream
	s.сhats[connectInfo.ChatID].mxStream.Unlock()
}

func (s *streamService) SendMessageToStreams(chatID string, msg *model.StreamMessage) error {
	for _, st := range s.сhats[chatID].streams {
		if err := st.Send(converter.ToGMessageFromModel(msg)); err != nil {
			return err
		}
	}
	return nil
}

func (s *streamService) DeleteStreamByUserID(chatID, userID string) {
	s.сhats[chatID].mxStream.Lock()
	delete(s.сhats[chatID].streams, userID)
	s.сhats[chatID].mxStream.Unlock()
}
