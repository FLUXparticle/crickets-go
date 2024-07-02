package service

import (
	"github.com/streadway/amqp"
	"log"
)

type MessageQueueProvider struct {
	logger       *log.Logger
	messageQueue MessageQueue
}

func NewMessageQueueProvider(logger *log.Logger) *MessageQueueProvider {
	return &MessageQueueProvider{
		logger: logger,
	}
}

func (p *MessageQueueProvider) GetMessageQueue() MessageQueue {
	if p.messageQueue == nil {
		mq, err := tryConnectRabbitMQ()
		if err != nil {
			p.logger.Println("WARNING:", err)
			mq = NewPubSub()
		}
		p.messageQueue = mq
	}
	return p.messageQueue
}

func tryConnectRabbitMQ() (MessageQueue, error) {
	// TODO Hostname aus Variable AMQP_HOST
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return NewRabbitMQ(channel), nil
}
