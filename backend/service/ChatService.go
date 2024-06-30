package service

import "crickets-go/repository"

type ChatService struct {
	pubSub *PubSub
}

func NewChatService(pubSub *PubSub) *ChatService {
	return &ChatService{pubSub: pubSub}
}

func (s *ChatService) ChatUpdates() chan *repository.Post {
	return s.pubSub.Subscribe("chat")
}

func (s *ChatService) SendChatMessage(post *repository.Post) {
	s.pubSub.Publish("chat", post)
}
