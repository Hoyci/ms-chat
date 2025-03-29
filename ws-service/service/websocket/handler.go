package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	coreTypes "github.com/hoyci/ms-chat/core/types"
	"github.com/hoyci/ms-chat/ws-service/service/rabbitmq"
	"github.com/hoyci/ms-chat/ws-service/types"
	"github.com/hoyci/ms-chat/ws-service/utils"
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

	clientID := uuid.New().String()

	fmt.Println("Headers:")
	for key, values := range r.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		conn.WriteJSON(types.WsErrorMessageResponse{
			ID:      clientID,
			Message: []string{"Missing X-User-ID on headers"},
			Status:  "missing_headers",
		})
		conn.Close()
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		conn.WriteJSON(types.WsErrorMessageResponse{
			ID:      clientID,
			Message: []string{"Invalid X-User-ID format"},
			Status:  "invalid_headers",
		})
		conn.Close()
		return
	}

	AddUserDeviceConnection(clientID, types.Connection{
		ClientID: clientID,
		UserID:   userID,
		Channel:  conn,
	})

	conn.WriteJSON(types.WsConnectionResponse{
		UserID:   userID,
		ClientID: clientID,
		Message:  "Successfully connected",
		Status:   "connected",
	})

	go manageConnection(conn, clientID)
}

func manageConnection(conn *websocket.Conn, clientID string) {
	defer func() {
		conn.Close()
		RemoveConnection(clientID)
		log.Printf("Connection %s closed", clientID)
	}()

	ch := rabbitmq.GetChannel()

	for {
		var msg coreTypes.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				log.Printf("Connection closed by client: %v", err)
				return
			}

			log.Printf("Read error: %v", err)
		}

		if err := utils.Validate.Struct(msg); err != nil {
			var errorMessages []string
			for _, e := range err.(validator.ValidationErrors) {
				errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' is invalid: %s", e.Field(), e.Tag()))
			}

			conn.WriteJSON(types.WsErrorMessageResponse{
				ID:      msg.ID,
				Message: errorMessages,
				Status:  "validation_error",
			})
			log.Printf("Client invalid message %s: %v", clientID, errorMessages)
			continue
		}

		msg.ID = uuid.New().String()
		msg.ClientID = clientID
		msg.CreatedAt = time.Now()
		msg.Status = "delivered"

		receiverDevices := GetUserDevicesConnections(msg.ReceiverID)

		connectionsCopy := make([]types.Connection, len(receiverDevices))
		copy(connectionsCopy, receiverDevices)

		if len(connectionsCopy) > 0 {
			for _, conn := range connectionsCopy {
				err := conn.Channel.WriteJSON(msg)
				if err != nil {
					log.Printf("An error occurred while sending the message to %s: %v", conn.ClientID, err)
					continue
				}
			}
		} else {
			msg.Status = "pending"
		}

		body, err := json.Marshal(msg)
		if err != nil {
			log.Printf("An unexpected error occurred while marshaling message: %v", err)
			return
		}

		retries := 0
		maxRetries := 3

		for retries < maxRetries {
			err = ch.Publish(
				"chat_events",
				"",
				false,
				false,
				amqp.Publishing{
					Headers:     amqp.Table{"persistence": "true"},
					ContentType: "application/json",
					Body:        body,
				},
			)
			if err == nil {
				log.Printf("Message %s published on exchange", msg.ID)
				break
			}
			log.Printf("Publish failed (attempt %d): %v", retries+1, err)
			time.Sleep(1 * time.Second)
			retries++
		}

		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("An unexpected error occurred while sending message to user: %v", err)
			continue
		}
	}
}
