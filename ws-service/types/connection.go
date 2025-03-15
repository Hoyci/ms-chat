package types

import "github.com/gorilla/websocket"

type Connection struct {
	ClientID string
	UserID   int
	Channel  *websocket.Conn
}

type WsConnectionResponse struct {
	UserID   int    `json:"user_id" validate:"required"`
	ClientID string `json:"client_id" validate:"required"`
	Message  string `json:"message" validate:"required"`
	Status   string `json:"status" validate:"required"`
}

type WsSuccessMessageResponse struct {
	TempID  string `json:"temp_id" validate:"required"`
	Message string `json:"message" validate:"required"`
	Status  string `json:"status" validate:"required"`
}

type WsErrorMessageResponse struct {
	TempID  string   `json:"temp_id" validate:"required"`
	Message []string `json:"message" validate:"required"`
	Status  string   `json:"status" validate:"required"`
}
