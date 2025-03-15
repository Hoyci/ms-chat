package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	coreTypes "github.com/hoyci/ms-chat/core/types"
	"github.com/hoyci/ms-chat/message-service/config"
	"github.com/hoyci/ms-chat/message-service/db"
	"github.com/hoyci/ms-chat/message-service/service/message"
	"github.com/hoyci/ms-chat/message-service/service/room"
	"github.com/hoyci/ms-chat/message-service/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MessageProcessor func(ctx context.Context, body []byte) error

var conn *amqp.Connection
var channel *amqp.Channel
var dbRepo *db.MongoRepository

func Init(repo *db.MongoRepository) {
	dbRepo = repo

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

func ConsumeQueue(channel *amqp.Channel, queueName string, processor MessageProcessor) {
	msgs, err := channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Printf("Failed to start consumption from %s: %v", queueName, err)
	}

	for msg := range msgs {
		err := processor(context.Background(), msg.Body)
		if err != nil {
			log.Printf("Error processing message: %v", err)
			msg.Nack(false, true)
		}
		msg.Ack(false)
	}
}

func ProcessChatMessage(ctx context.Context, msgBody []byte) error {
	var wsMessage coreTypes.Message
	json.Unmarshal(msgBody, &wsMessage)

	b, err := json.Marshal(wsMessage)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(b))

	// Pegar ou cria o canal da mensagem
	roomStore := room.GetRoomStore(dbRepo)
	messageStore := message.GetMessageStore(dbRepo)
	log.Printf("Pegando ou criando canal da mensagem: %+v", wsMessage)
	room, err := roomStore.GetOrCreate(context.Background(), wsMessage.RoomID)
	if err != nil {
		log.Printf("Error with GetOrCreateRoom: %v", err)
		return err
	}

	// Persistir mensagem no banco de dados com status pending
	log.Printf("Persistindo mensagem2: %+v", wsMessage)
	messageID, err := messageStore.Create(
		context.Background(),
		types.Message{
			ID:         bson.NewObjectID(),
			RoomID:     room.ID.Hex(),
			SenderID:   wsMessage.SenderID,
			ReceiverID: wsMessage.ReceiverID,
			Content:    wsMessage.Content,
			Status:     types.StatusPending,
			CreatedAt:  time.Now(),
			UpdatedAt:  nil,
			DeletedAt:  nil,
		},
	)
	if err != nil {
		log.Printf("Error persisting message: %v", err)
		return err
	}

	log.Printf("message: %s", messageID)

	// Verifica se o usuário está online
	// Se estiver, atualizada status da mensagem para delivered e publica na fila de "broadcast"
	return nil
}
