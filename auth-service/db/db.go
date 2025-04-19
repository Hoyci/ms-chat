package db

import (
	"database/sql"
	"log"

	"github.com/hoyci/ms-chat/auth-service/config"
	_ "github.com/lib/pq"
)

func NewPGStorage() *sql.DB {
	db, err := sql.Open("postgres", config.Envs.DatabaseURL.String())
	if err != nil {
		log.Fatalf("Failed to start connection")
	}

	return db
}
