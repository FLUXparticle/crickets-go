package service

import (
	"crickets-go/config"
	"crickets-go/data"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQ struct {
	logger          *log.Logger
	config          *config.Config
	rabbitMQChannel *amqp.Channel
}

func NewRabbitMQ(logger *log.Logger, config *config.Config, rabbitMQChannel *amqp.Channel) *RabbitMQ {
	return &RabbitMQ{
		logger:          logger,
		config:          config,
		rabbitMQChannel: rabbitMQChannel,
	}
}

func (r *RabbitMQ) Subscribe(topic string) chan *data.Post {
	// Exchange erzeugen, falls nötig
	exchangeName := topic + ".exchange"
	err := r.rabbitMQChannel.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		r.logger.Panicf("Failed to declare an exchange: %s", err)
	}

	// Erzeuge eine Queue für den Client
	queueName := topic
	if len(r.config.Hostname) > 0 {
		queueName += "_" + r.config.Hostname
	}
	queueName += ".queue"
	_, err = r.rabbitMQChannel.QueueDeclare(
		queueName, // name
		false,     // durable
		true,      // delete when unused
		true,      // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		r.logger.Panicf("Failed to declare a queue: %s", err)
	}

	// Binde die Queue an den Fanout-Exchange
	err = r.rabbitMQChannel.QueueBind(
		queueName,    // queue name
		"",           // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		r.logger.Panicf("Failed to bind a queue: %s", err)
	}

	messages, err := r.rabbitMQChannel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		// Handle error
		panic(err)
	}

	ch := make(chan *data.Post)
	go func() {
		for msg := range messages {
			var post data.Post
			err := json.Unmarshal(msg.Body, &post)
			if err != nil {
				r.logger.Println("Error receiving:", err.Error())
				continue
			}

			// Weiterleitung an den Client
			ch <- &post
		}
	}()

	return ch
}

func (r *RabbitMQ) Unsubscribe(topic string, ch chan *data.Post) {
	panic("implement me")
}

func (r *RabbitMQ) Publish(topic string, post *data.Post) error {
	// Marshalling des Posts in JSON
	body, err := json.Marshal(post)
	if err != nil {
		return err
	}

	// Nachricht veröffentlichen
	exchangeName := topic + ".exchange"
	err = r.rabbitMQChannel.Publish(
		exchangeName, // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}
