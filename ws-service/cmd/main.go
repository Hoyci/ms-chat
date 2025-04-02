package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"

	"github.com/hoyci/ms-chat/ws-service/config"
	"github.com/hoyci/ms-chat/ws-service/service/rabbitmq"
	"github.com/hoyci/ms-chat/ws-service/service/websocket"
	"github.com/hoyci/ms-chat/ws-service/utils"
)

func main() {
	path := fmt.Sprintf("0.0.0.0:%d", config.Envs.Port)

	rabbitmq.Init()
	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {

		}
	}(rabbitmq.GetChannel())

	utils.InitValidator()

	websocket.RegisterRoutes()

	log.Println("Listening on:", path)
	err := http.ListenAndServe(path, nil)
	if err != nil {
		return
	}
}
