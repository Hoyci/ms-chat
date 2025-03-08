package websocket

import "net/http"

func RegisterRoutes() {
	http.HandleFunc("/ws", HandleWebsocket)
}
