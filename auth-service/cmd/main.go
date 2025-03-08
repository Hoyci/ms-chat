package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hoyci/ms-chat/auth-service/cmd/api"
	"github.com/hoyci/ms-chat/auth-service/config"
	"github.com/hoyci/ms-chat/auth-service/db"
	"github.com/hoyci/ms-chat/auth-service/service/auth"
	"github.com/hoyci/ms-chat/auth-service/service/healthcheck"
	"github.com/hoyci/ms-chat/auth-service/service/user"
	"github.com/hoyci/ms-chat/auth-service/utils"
)

// @title Auth Service API dwada
// @version 1.0
// @description API para gestão de usuário e autenticação
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	db := db.NewPGStorage()
	path := fmt.Sprintf("0.0.0.0:%d", config.Envs.Port)

	apiServer := api.NewApiServer(path, db)

	healthCheckHandler := healthcheck.NewHealthCheckHandler(config.Envs)

	userStore := user.NewUserStore(db)
	userHandler := user.NewUserHandler(userStore)

	authStore := auth.NewAuthStore(db)
	uuidGen := &utils.UUIDGeneratorUtil{}
	authHandler := auth.NewAuthHandler(userStore, authStore, uuidGen)

	apiServer.SetupRouter(healthCheckHandler, userHandler, authHandler)

	log.Println("Listening on:", path)
	http.ListenAndServe(path, apiServer.Router)
}
