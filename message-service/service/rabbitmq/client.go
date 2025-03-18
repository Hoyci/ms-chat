package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	coreTypes "github.com/hoyci/ms-chat/core/types"
	"github.com/hoyci/ms-chat/message-service/config"
	"github.com/hoyci/ms-chat/message-service/db"
	"github.com/hoyci/ms-chat/message-service/service/message"
	"github.com/hoyci/ms-chat/message-service/service/room"
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

	roomStore := room.GetRoomStore(dbRepo)
	messageStore := message.GetMessageStore(dbRepo)

	room, err := roomStore.GetOrCreate(context.Background(), []int{wsMessage.SenderID, wsMessage.ReceiverID})
	if err != nil {
		log.Printf("Error with GetOrCreateRoom: %v", err)
		return err
	}

	messageID, err := messageStore.Create(
		context.Background(),
		map[string]any{
			"_id":         bson.NewObjectID(),
			"room_id":     room.ID,
			"sender_id":   wsMessage.SenderID,
			"receiver_id": wsMessage.ReceiverID,
			"content":     wsMessage.Content,
			"status":      wsMessage.Status,
			"created_at":  wsMessage.CreatedAt,
			"updated_at":  nil,
			"deleted_at":  nil,
		},
	)
	if err != nil {
		log.Printf("Error persisting message: %v", err)
		return err
	}

	log.Printf("message: %s", messageID.Hex())

	return nil
}
