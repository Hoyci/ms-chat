package config

import (
	"log"
	"os"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Port                   int    `env:"PORT" envDefault:"8080"`
	Environment            string `env:"ENVIRONMENT" envDefault:"development"`
	DatabaseURL            string `env:"DATABASE_URL" envDefault:"postgres://user:password@postgres:5432/postgres?sslmode=disable"`
	KeysPath               string `env:"KEYS_PATH" envDefault:"./keys"`
	PublicKeyFilename      string `env:"PUBLIC_KEY_FILENAME" envDefault:"public_key_access.pem"`
	TestPrivateKeyFilename string `env:"TEST_PRIVATE_KEY_FILENAME" envDefault:"test_private_key.pem"`
	TestPublicKeyFilename  string `env:"TEST_PUBLIC_KEY_FILENAME" envDefault:"test_public_key.pem"`
}

var Envs = initConfig()

func initConfig() Config {
	if err := godotenv.Load("C:\\Users\\Administrador\\golang\\ms-chat\\contacts-service\\.env"); err != nil {
		if os.IsNotExist(err) {
			log.Println(".env file not found, using default environment variables")
		} else {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	return cfg
}
