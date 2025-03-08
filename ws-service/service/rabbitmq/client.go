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

	err = channel.ExchangeDeclare(
		"chat_events",
		"headers",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	createQueue("persistence_queue", amqp.Table{
		"x-match":     "any",
		"persistence": "true",
	})

	createQueue("broadcast_queue", amqp.Table{
		"x-match":   "any",
		"broadcast": "true",
	})
}

func createQueue(name string, headers amqp.Table) {
	q, _ := channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)

	channel.QueueBind(
		q.Name,
		"",
		"chat_events",
		false,
		headers,
	)
}

func GetChannel() *amqp.Channel {
	return channel
}
