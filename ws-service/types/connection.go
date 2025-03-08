package types

import "github.com/gorilla/websocket"

type Connection struct {
	ClientID string
	UserID   int
	Rooms    map[string]bool
	Channel  *websocket.Conn
}
