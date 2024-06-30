package service

import (
	"crickets-go/repository"
	"sync"
)

type PubSub struct {
	channelsMap map[string][]chan *repository.Post
	mu          sync.RWMutex
}

func NewPubSub() *PubSub {
	return &PubSub{
		channelsMap: make(map[string][]chan *repository.Post),
	}
}

func (ps *PubSub) Subscribe(topic string) chan *repository.Post {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan *repository.Post)
	ps.channelsMap[topic] = append(ps.channelsMap[topic], ch)
	return ch
}

func (ps *PubSub) Unsubscribe(topic string, ch chan *repository.Post) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	channels := ps.channelsMap[topic]
	for i, c := range channels {
		if c == ch {
			ps.channelsMap[topic] = append(channels[:i], channels[i+1:]...)
			close(c)
			break
		}
	}
}

func (ps *PubSub) Publish(topic string, post *repository.Post) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if channels, found := ps.channelsMap[topic]; found {
		for _, ch := range channels {
			ch <- post
		}
	}
}
