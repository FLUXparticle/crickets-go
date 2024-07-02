package service

import (
	"crickets-go/repository"
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

func (s *ChatService) ChatUpdates() chan *repository.Post {
	mq := s.messageQueueProvider.GetMessageQueue()
	return mq.Subscribe(chatTopic)
}

func (s *ChatService) SendChatMessage(post *repository.Post) error {
	mq := s.messageQueueProvider.GetMessageQueue()
	return mq.Publish(chatTopic, post)
}
