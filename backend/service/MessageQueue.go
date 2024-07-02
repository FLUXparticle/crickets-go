package service

import "crickets-go/repository"

type MessageQueue interface {
	Subscribe(topic string) chan *repository.Post
	Unsubscribe(topic string, ch chan *repository.Post)
	Publish(topic string, post *repository.Post) error
}
