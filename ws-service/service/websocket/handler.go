package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hoyci/ms-chat/ws-service/service/rabbitmq"
	"github.com/hoyci/ms-chat/ws-service/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error updating connection to websocket:", err)
		return
	}

	userIDStr := r.URL.Query().Get("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Println("Invalid userId:", err)
		return
	}

	clientID := uuid.New().String()

	AddConnection(clientID, types.Connection{
		ClientID: clientID,
		UserID:   userID,
		Rooms:    make(map[string]bool),
		Channel:  conn,
	})

	go manageConnection(conn, clientID)
}

func manageConnection(conn *websocket.Conn, clientID string) {
	defer func() {
		RemoveConnection(clientID)
		conn.Close()
	}()

	ch := rabbitmq.GetChannel()

	for {
		var msg types.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("Unexpected error:", err)
			} else {
				log.Println("Connection closed by client:", err)
			}
			return
		}

		switch msg.Type {
		case "join_room":
			mu.Lock()
			if connData, exists := connections[clientID]; exists {
				connData.Rooms[msg.Room] = true
				connections[clientID] = connData
			}
			mu.Unlock()

		case "leave_room":
			mu.Lock()
			if connData, exists := connections[clientID]; exists {
				delete(connData.Rooms, msg.Room)
				connections[clientID] = connData
			}
			mu.Unlock()
		}

		msg.ClientID = clientID
		msg.Timestamp = time.Now().UTC().Format(time.RFC3339)

		body, _ := json.Marshal(msg)
		log.Printf("[PUBLISH] Enviando mensagem para exchange. Body: %s", string(body))

		err = ch.Publish(
			"chat_events",
			"",
			false,
			false,
			amqp.Publishing{
				Headers: amqp.Table{
					"persistence": "true",
					"broadcast":   "true",
				},
				ContentType: "application/json",
				Body:        body,
			},
		)

		if err != nil {
			log.Println("Failed to publish message:", err)
		} else {
			log.Println("[SUCESSO] Mensagem publicada no exchange")
		}
	}
}
