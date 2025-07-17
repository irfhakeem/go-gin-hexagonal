package ports

import "github.com/rabbitmq/amqp091-go"

type MessageQueueManager interface {
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args map[string]any) (amqp091.Queue, error)
	Publisher(exchanged, name string, mandatory, immediate bool, msg []byte) error
	Consumer(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args map[string]any) (<-chan amqp091.Delivery, error)
}
