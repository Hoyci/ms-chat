package websocket

import (
	"encoding/json"
	"log"

	"github.com/hoyci/ms-chat/ws-service/service/rabbitmq"
	"github.com/hoyci/ms-chat/ws-service/types"
)

func StartBroadcastConsumer() {
	ch := rabbitmq.GetChannel()

	msgs, _ := ch.Consume(
		"broadcast_queue",
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	for msg := range msgs {
		var message types.Message

		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Println("Error unmarshaling message", err)
			msg.Nack(false, true)
			continue
		}

		userDevices := GetUserDevicesConnections(message.UserID)
		deliveryFailed := false

		for _, conn := range userDevices {
			// if ok := conn.UserID == message.UserID; !ok {
			// 	continue
			// }

			// if conn.ClientID == message.ClientID {
			// 	continue
			// }

			if err := conn.Channel.WriteJSON(message); err != nil {
				log.Printf("An error occurred while sending the message to %s: %v", conn.ClientID, err)
				deliveryFailed = true
			}
		}

		if deliveryFailed {
			msg.Nack(false, true)
		} else {
			msg.Ack(false)
		}
	}
}
