package service

import "crickets-go/data"

type MessageQueue interface {
	Subscribe(topic string) chan *data.Post
	Unsubscribe(topic string, ch chan *data.Post)
	Publish(topic string, post *data.Post) error
}
