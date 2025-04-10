package main

import (
	"fmt"
	"github.com/hoyci/ms-chat/auth-service/keys"
	"github.com/hoyci/ms-chat/auth-service/service/crypt"
	coreUtils "github.com/hoyci/ms-chat/core/utils"
	"log"
	"net/http"

	"github.com/hoyci/ms-chat/auth-service/cmd/api"
	"github.com/hoyci/ms-chat/auth-service/config"
	"github.com/hoyci/ms-chat/auth-service/db"
	"github.com/hoyci/ms-chat/auth-service/service/auth"
	"github.com/hoyci/ms-chat/auth-service/service/healthcheck"
	"github.com/hoyci/ms-chat/auth-service/service/user"
)

// @title Auth Service API
// @version 1.0
// @description API para gestão de usuário e autenticação
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	keys.LoadRunKeys()
	pgStorage := db.NewPGStorage()
	path := fmt.Sprintf("0.0.0.0:%d", config.Envs.Port)

	apiServer := api.NewServer(path, pgStorage)

	healthCheckHandler := healthcheck.NewHealthCheckHandler(config.Envs)

	userStore := user.NewUserStore(pgStorage)
	passwordStore := &crypt.BcryptPasswordStore{}
	passwordHandler := crypt.PasswordHandler(passwordStore)
	userHandler := user.NewUserHandler(userStore, passwordHandler)

	authStore := auth.NewAuthStore(pgStorage)
	uuidGen := &coreUtils.UUIDGeneratorUtil{}
	authHandler := auth.NewAuthHandler(userStore, authStore, uuidGen, passwordHandler)

	apiServer.SetupRouter(healthCheckHandler, userHandler, authHandler)

	log.Println("Listening on:", path)
	err := http.ListenAndServe(path, apiServer.Router)
	if err != nil {
		log.Panic("Failed to start server: " + err.Error())
	}
}
