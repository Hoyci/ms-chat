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

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Println("Invalid userID:", err)
		return
	}

	clientID := uuid.New().String()

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
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("Connection closed by client: %v", err)
				return
			}

			if err := utils.Validate.Struct(msg); err != nil {
				var errorMessages []string
				for _, e := range err.(validator.ValidationErrors) {
					errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' is invalid: %s", e.Field(), e.Tag()))
				}

				conn.WriteJSON(types.WsErrorMessageResponse{
					TempID:  msg.TempID,
					Message: errorMessages,
					Status:  "validation_error",
				})
				log.Printf("Client invalid message %s: %v", clientID, errorMessages)
				continue
			}
		}

		msg.TempID = uuid.New().String()
		msg.ClientID = clientID
		msg.Timestamp = time.Now().UTC().Format(time.RFC3339)

		body, err := json.Marshal(msg)
		if err != nil {
			log.Printf("An unexpected error occurred while marshaling message: %v", err)
			return
		}

		receiverDevices := GetUserDevicesConnections(msg.ReceiverID)

		for _, conn := range receiverDevices {
			if err := conn.Channel.WriteJSON(msg); err != nil {
				log.Printf("An error occurred while sending the message to %s: %v", conn.ClientID, err)
			}
		}

		err = ch.Publish(
			"chat_events",
			"",
			false,
			false,
			amqp.Publishing{
				Headers: amqp.Table{
					"persistence": "true",
				},
				ContentType: "application/json",
				Body:        body,
			},
		)

		if err != nil {
			log.Println("Failed to publish message:", err)
		}

		if err := conn.WriteJSON(types.WsSuccessMessageResponse{
			TempID:  msg.TempID,
			Status:  "sent",
			Message: "Message sent successfully",
		}); err != nil {
			log.Printf("An unexpected error occurred while sending message to user: %v", err)
			return
		}
		log.Println("Message published on exchange")
	}
}
