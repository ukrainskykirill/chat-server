package chats

import (
	"github.com/ukrainskykirill/chat-server/internal/model"
)

func (s *chatServ) ConnectChat(connectInfo *model.ConnectInfo) error {
	chatChan, err := s.streamService.GetChatChannel(connectInfo.Stream.Context(), connectInfo.ChatID)
	if err != nil {
		return err
	}

	s.streamService.AddStreamByUserID(connectInfo)

	for {
		select {
		case msg, okCh := <-chatChan:
			if !okCh {
				return nil
			}
			s.streamService.SendMessageToStreams(connectInfo.ChatID, msg)
		case <-connectInfo.Stream.Context().Done():
			s.streamService.DeleteStreamByUserID(connectInfo.ChatID, connectInfo.UserID)
			return nil
		}

	}

	return nil
}
