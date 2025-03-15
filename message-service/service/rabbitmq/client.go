package rabbitmq

import (
	"encoding/json"
	"log"

	coreTypes "github.com/hoyci/ms-chat/core/types"
	"github.com/hoyci/ms-chat/message-service/config"
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
}

func GetChannel() *amqp.Channel {
	return channel
}

func ConsumeQueue(queue string, channel *amqp.Channel) {
	msgs, err := channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Printf("Failed to start consumption from %s: %v", queue, err)
	}

	for msg := range msgs {
		var message coreTypes.Message
		json.Unmarshal(msg.Body, &message)

		// Pegar ou cria o canal da mensagem
		log.Printf("Pegando ou criando canal da mensagem: %+v", message)
		// Persistir mensagem no banco de dados com status pending
		log.Printf("Persistindo mensagem: %+v", message)
		// Verifica se o usuário está online
		// Se estiver, atualizada status da mensagem para delivered e publica na fila de "broadcast"

		msg.Ack(false)
	}
}
