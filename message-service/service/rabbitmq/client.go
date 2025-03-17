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
	"github.com/hoyci/ms-chat/message-service/service/redis"
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

	// Pegar ou cria o canal da mensagem
	roomStore := room.GetRoomStore(dbRepo)
	messageStore := message.GetMessageStore(dbRepo)
	room, err := roomStore.GetOrCreate(context.Background(), wsMessage.RoomID, []int{wsMessage.SenderID, wsMessage.ReceiverID})
	if err != nil {
		log.Printf("Error with GetOrCreateRoom: %v", err)
		return err
	}

	// Persistir mensagem no banco de dados com status pendin
	messageID, err := messageStore.Create(
		context.Background(),
		types.Message{
			ID:        bson.NewObjectID(),
			RoomID:    room.ID,
			UserID:    wsMessage.SenderID,
			Content:   wsMessage.Content,
			Status:    types.StatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: nil,
			DeletedAt: nil,
		},
	)
	if err != nil {
		log.Printf("Error persisting message: %v", err)
		return err
	}

	log.Printf("message: %s", messageID.Hex())

	// Verifica se o usuário está online
	if isOnline := redis.IsUserOnline(wsMessage.ReceiverID); isOnline {
		// atualiza status da mensagem para delivered

		// publica na fila de "broadcast"
		body, err := json.Marshal(types.Message{
			ID:        messageID,
			RoomID:    room.ID,
			UserID:    wsMessage.ReceiverID,
			Content:   "teste",
			Status:    "delivered",
			CreatedAt: time.Now(),
			UpdatedAt: nil,
			DeletedAt: nil,
		})
		if err != nil {
			log.Printf("An unexpected error occurred while marshaling message: %v", err)
			return nil
		}

		ch := GetChannel()
		err = ch.Publish(
			"chat_events",
			"",
			false,
			false,
			amqp.Publishing{
				Headers: amqp.Table{
					"broadcast": "true",
				},
				ContentType: "application/json",
				Body:        body,
			},
		)

		if err != nil {
			log.Println("Failed to publish message:", err)
			return nil
		}
	}

	return nil
}
