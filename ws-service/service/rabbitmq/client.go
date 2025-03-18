package rabbitmq

import (
	"encoding/json"
	"log"
	"time"

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

	createQueue("persistence_queue", amqp.Table{"x-match": "any", "persistence": "true"})
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

func createQueue(name string, headers amqp.Table) {
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

func PublishUserOnlineEvent(userID int) {
	body, err := json.Marshal(map[string]interface{}{
		"user_id": userID,
		"event":   "user_online",
		"time":    time.Now().UTC(),
	})

	if err != nil {
		log.Printf("Failed to marshal event body %v", err)
	}

	err = channel.Publish(
		"user_events",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Printf("Failed to publish user online event: %v", err)
	}
}
