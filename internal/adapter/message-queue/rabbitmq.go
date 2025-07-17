package mq

import (
	"go-gin-hexagonal/internal/domain/ports"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	ch *amqp091.Channel
}

func NewRabbitMQ(ch *amqp091.Channel) ports.MessageQueueManager {
	return &RabbitMQ{
		ch: ch,
	}
}

func (r *RabbitMQ) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args map[string]any) (amqp091.Queue, error) {
	return r.ch.QueueDeclare(
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args,
	)
}

func (r *RabbitMQ) Publisher(exchanged, name string, mandatory, immediate bool, msg []byte) error {
	return r.ch.Publish(
		exchanged, // exchange
		name,      // routing key
		mandatory, // mandatory
		immediate, // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
}

func (r *RabbitMQ) Consumer(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args map[string]any) (<-chan amqp091.Delivery, error) {
	msgs, err := r.ch.Consume(
		queue,
		consumer,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		args,
	)

	if err != nil {
		return nil, err
	}

	return msgs, nil
}
