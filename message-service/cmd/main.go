package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hoyci/ms-chat/message-service/cmd/api"
	"github.com/hoyci/ms-chat/message-service/config"
	"github.com/hoyci/ms-chat/message-service/db"
	"github.com/hoyci/ms-chat/message-service/service/healthcheck"
	"github.com/hoyci/ms-chat/message-service/service/room"
)

// @title Message Service API
// @version 1.0
// @description API para gestão de conversas e mensagens
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	dbRepo := db.NewMongoRepository(config.Envs)
	path := fmt.Sprintf("0.0.0.0:%d", config.Envs.Port)

	apiServer := api.NewApiServer(path)

	healthCheckHandler := healthcheck.NewHealthCheckHandler(config.Envs)

	roomStore := room.NewRoomStore(dbRepo)
	roomHandler := room.NewRoomHandler(roomStore)

	apiServer.SetupRouter(
		healthCheckHandler,
		roomHandler,
	)

	log.Println("Listening on:", path)
	http.ListenAndServe(path, apiServer.Router)
}
