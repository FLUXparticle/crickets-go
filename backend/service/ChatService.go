package service

import (
	"crickets-go/data"
)

const chatTopic = "chat_queue"

type ChatService struct {
	messageQueueProvider *MessageQueueProvider
}

func NewChatService(messageQueueProvider *MessageQueueProvider) *ChatService {
	return &ChatService{
		messageQueueProvider: messageQueueProvider,
	}
}

func (s *ChatService) ChatUpdates() chan *data.Post {
	mq := s.messageQueueProvider.GetMessageQueue()
	return mq.Subscribe(chatTopic)
}

func (s *ChatService) SendChatMessage(post *data.Post) error {
	mq := s.messageQueueProvider.GetMessageQueue()
	return mq.Publish(chatTopic, post)
}
