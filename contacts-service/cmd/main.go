package main

import (
	"fmt"
	"github.com/hoyci/ms-chat/contacts-service/cmd/api"
	"github.com/hoyci/ms-chat/contacts-service/config"
	"github.com/hoyci/ms-chat/contacts-service/db"
	"github.com/hoyci/ms-chat/contacts-service/services/contacts"
	"log"
	"net/http"
)

func main() {
	pgStorage := db.NewPGStorage()
	path := fmt.Sprintf("0.0.0.0:%d", config.Envs.Port)

	apiServer := api.NewServer(path, pgStorage)

	contactStore := contacts.NewContactStore(pgStorage)
	contactHandler := contacts.NewContactHandler(contactStore)

	apiServer.SetupRouter(contactHandler)

	log.Println("Listening on:", path)
	err := http.ListenAndServe(path, apiServer.Router)
	if err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
