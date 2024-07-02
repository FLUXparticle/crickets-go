package service

import (
	"crickets-go/config"
	"github.com/streadway/amqp"
	"log"
)

type MessageQueueProvider struct {
	logger       *log.Logger
	config       *config.Config
	messageQueue MessageQueue
}

func NewMessageQueueProvider(logger *log.Logger, config *config.Config) *MessageQueueProvider {
	return &MessageQueueProvider{
		logger: logger,
		config: config,
	}
}

func (p *MessageQueueProvider) GetMessageQueue() MessageQueue {
	if p.messageQueue == nil {
		mq, err := p.tryConnectRabbitMQ()
		if err != nil {
			p.logger.Println("WARNING:", err)
			mq = NewPubSub()
		}
		p.messageQueue = mq
	}
	return p.messageQueue
}

func (p *MessageQueueProvider) tryConnectRabbitMQ() (MessageQueue, error) {
	host := p.config.AmqpHost
	conn, err := amqp.Dial("amqp://guest:guest@" + host + ":5672/")
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return NewRabbitMQ(p.logger, p.config, channel), nil
}
