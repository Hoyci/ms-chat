package rabbitmq

import (
	"log"

	"github.com/hoyci/ms-chat/ws-service/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var channel *amqp.Channel

func Init() {
	var err error
	conn, err = amqp.Dial(config.Envs.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	channel, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}

	createExchange("chat_events", "headers")
	createExchange("user_events", "fanout")

	persistenceQueue := createQueue(config.Envs.PersistenceQueueName)

	channel.QueueBind(
		persistenceQueue.Name,
		"",
		"chat_events",
		false,
		amqp.Table{"x-match": "any", "persistence": "true"},
	)
}

func createExchange(name string, kind string) {
	err := channel.ExchangeDeclare(
		name,
		kind,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange %s: %v", name, err)
	}
}

func createQueue(name string) amqp.Queue {
	q, err := channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Printf("Failed to declare queue %s", name)
	}

	return q
}

func GetChannel() *amqp.Channel {
	return channel
}
