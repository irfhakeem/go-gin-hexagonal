package mq

import (
	"go-gin-hexagonal/pkg/config"

	"github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQConnection(cfg *config.RabbitMQConfig) *amqp091.Channel {
	dsn := cfg.DSN()
	conn, err := amqp091.Dial(dsn)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return ch
}
