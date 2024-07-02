package service

import (
	"crickets-go/repository"
	"encoding/json"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	rabbitMQChannel *amqp.Channel
}

func NewRabbitMQ(rabbitMQChannel *amqp.Channel) *RabbitMQ {
	return &RabbitMQ{rabbitMQChannel: rabbitMQChannel}
}

func (r *RabbitMQ) Subscribe(topic string) chan *repository.Post {
	messages, err := r.rabbitMQChannel.Consume(
		topic, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		// Handle error
		panic(err)
	}

	ch := make(chan *repository.Post)
	go func() {
		for msg := range messages {
			var post repository.Post
			err := json.Unmarshal(msg.Body, &post)
			if err != nil {
				// Handle error
				continue
			}

			//fmt.Println("Received:", post.Content)

			// Weiterleitung an den Client
			ch <- &post
		}
	}()

	return ch
}

func (r *RabbitMQ) Unsubscribe(topic string, ch chan *repository.Post) {
	panic("implement me")
}

func (r *RabbitMQ) Publish(topic string, post *repository.Post) error {
	// Marshalling des Posts in JSON
	body, err := json.Marshal(post)
	if err != nil {
		return err
	}

	// Nachricht veröffentlichen
	err = r.rabbitMQChannel.Publish(
		"",    // exchange
		topic, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}
