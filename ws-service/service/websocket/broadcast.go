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

	log.Println("[CONSUMER] Iniciando consumer do broadcast...")
	for msg := range msgs {
		log.Printf("[BROADCAST] Mensagem recebida: %s", string(msg.Body))
		var message types.Message
		json.Unmarshal(msg.Body, &message)

		connections := GetRoomConnections(message.Room)
		for _, conn := range connections {
			log.Printf("Verificando conexão %s para sala %s", conn.ClientID, message.Room)
			if _, ok := conn.Rooms[message.Room]; ok {
				log.Printf("Conexão %s está na sala", conn.ClientID)
				if conn.ClientID != message.ClientID {
					log.Printf("Enviando para %s", conn.ClientID)
					err := conn.Channel.WriteJSON(message)
					if err != nil {
						log.Printf("Erro ao enviar para %s: %v", conn.ClientID, err)
					}
				}
			}
		}

		msg.Ack(false)
	}
}
