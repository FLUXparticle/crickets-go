package service

import (
	"crickets-go/repository"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

type ChatService struct {
	pubSub          *PubSub
	rabbitMQConn    *amqp.Connection
	rabbitMQChannel *amqp.Channel
}

func NewRabbitConnection() *amqp.Connection {
	// TODO Hostname aus Variable AMQP_HOST
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		// TODO panic(err)
		return nil
	}
	return conn
}

func NewChatService(pubSub *PubSub, rabbitMQConn *amqp.Connection) *ChatService {
	// TODO rabbitMQChannel, err := rabbitMQConn.Channel()
	//if err != nil {
	// TODO panic(err)
	//}

	return &ChatService{
		pubSub:       pubSub,
		rabbitMQConn: rabbitMQConn,
		// TODO rabbitMQChannel: rabbitMQChannel,
	}
}

func (s *ChatService) ChatUpdates() chan *repository.Post {
	messages, err := s.rabbitMQChannel.Consume(
		"chat_queue", // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
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

	//return s.pubSub.Subscribe("chat")
}

func (s *ChatService) SendChatMessage(post *repository.Post) {
	// Marshalling des Posts in JSON
	body, err := json.Marshal(post)
	if err != nil {
		// Handle error
		panic(err)
	}

	// Nachricht veröffentlichen
	err = s.rabbitMQChannel.Publish(
		"",           // exchange
		"chat_queue", // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		// Handle error
		fmt.Println("Error:", err.Error())
		panic(err)
	}

	//fmt.Println("Sent:", post.Content)

	//s.pubSub.Publish("chat", post)
}
