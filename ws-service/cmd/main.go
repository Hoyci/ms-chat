package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hoyci/ms-chat/ws-service/config"
	"github.com/hoyci/ms-chat/ws-service/service/rabbitmq"
	"github.com/hoyci/ms-chat/ws-service/service/websocket"
	"github.com/hoyci/ms-chat/ws-service/utils"
)

func main() {
	path := fmt.Sprintf("localhost:%d", config.Envs.Port)

	rabbitmq.Init()
	defer rabbitmq.GetChannel().Close()
	// websocket.StartBroadcastConsumer()

	utils.InitValidator()

	websocket.RegisterRoutes()

	log.Println("Listening on:", path)
	http.ListenAndServe(path, nil)
}
