package websocket

import (
	"encoding/json"

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
		var broadcast types.BroadcastMessage
		if err := json.Unmarshal(msg.Body, &broadcast); err != nil {
			continue
		}

		receiverDevices := GetUserDevicesConnections(broadcast.UserID)

		connectionsCopy := make([]types.Connection, len(receiverDevices))
		copy(connectionsCopy, receiverDevices)

		if connections := GetUserDevicesConnections(broadcast.UserID); len(connections) > 0 {
			for _, conn := range connections {
				conn.Channel.WriteJSON(broadcast.Messages)
			}
		}

		msg.Ack(false)
	}
}
