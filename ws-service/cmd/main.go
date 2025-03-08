package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hoyci/ms-chat/ws-service/config"
	"github.com/hoyci/ms-chat/ws-service/service/rabbitmq"
	"github.com/hoyci/ms-chat/ws-service/service/websocket"
)

func main() {
	path := fmt.Sprintf("0.0.0.0:%d", config.Envs.Port)

	rabbitmq.Init()
	defer rabbitmq.GetChannel().Close()
	go websocket.StartBroadcastConsumer()

	websocket.RegisterRoutes()

	log.Println("Listening on:", path)
	http.ListenAndServe(path, nil)
}
