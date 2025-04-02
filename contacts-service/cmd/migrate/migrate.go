package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hoyci/ms-chat/contacts-service/config"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", config.Envs.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}(db)

	driver, err := pgMigrate.WithInstance(db, &pgMigrate.Config{})
	if err != nil {
		log.Fatalf("Error creating migrate driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"public",
		driver,
	)
	if err != nil {
		log.Fatalf("Error creating migrate instance: %v", err)
	}

	cmd := os.Args[len(os.Args)-1]
	switch cmd {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migrations applied successfully.")
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migrations reverted successfully.")
	default:
		log.Println("No command provided. Use 'up' or 'down'.")
	}
}
