package service

import (
	"crickets-go/config"
	"crickets-go/data"
	"time"
)

const chatTopic = "chat"

type ChatService struct {
	config               *config.Config
	messageQueueProvider *MessageQueueProvider
}

func NewChatService(config *config.Config, messageQueueProvider *MessageQueueProvider) *ChatService {
	return &ChatService{
		config:               config,
		messageQueueProvider: messageQueueProvider,
	}
}

func (s *ChatService) ChatUpdates() chan *data.Post {
	mq := s.messageQueueProvider.GetMessageQueue()
	return mq.Subscribe(chatTopic)
}

func (s *ChatService) SendChatMessage(user *data.User, content string) error {
	post := &data.Post{
		Creator: &data.User{
			Server:   s.config.Hostname,
			Username: user.Username,
		},
		Content:   content,
		CreatedAt: time.Now(),
	}
	mq := s.messageQueueProvider.GetMessageQueue()
	return mq.Publish(chatTopic, post)
}
