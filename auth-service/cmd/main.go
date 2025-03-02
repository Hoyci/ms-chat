package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/hoyci/auth-service/cmd/api"
	"github.com/hoyci/auth-service/types"
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
	var cfg types.Config
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	path := fmt.Sprintf("0.0.0.0:%d", cfg.Port)

	apiServer := api.NewApiServer(path)
	apiServer.SetupRouter()

	log.Println("Listening on:", path)
	http.ListenAndServe(path, apiServer.Router)
}
