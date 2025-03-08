package config

import (
	"log"
	"os"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Port                          int    `env:"PORT" envDefault:"8080"`
	Environment                   string `env:"ENVIRONMENT" envDefault:"development"`
	DatabaseURL                   string `env:"DATABASE_URL" envDefault:"postgres://user:password@postgres:5432/postgres?sslmode=disable"`
	AccessJWTSecret               string `env:"ACCESS_JWT_SECRET" envDefault:"UM_ACCESS_TOKEN_MTO_DIFICIL"`
	AccessJWTExpirationInSeconds  int    `env:"ACCESS_JWT_EXPIRATION" envDefault:"3600"`
	RefreshJWTSecret              string `env:"REFRESH_JWT_SECRET" envDefault:"UM_REFRESH_TOKEN_MTO_DIFICIL"`
	RefreshJWTExpirationInSeconds int    `env:"REFRESH_JWT_EXPIRATION" envDefault:"604800"`
	RootPath                      string `env:"ROOT_PATH" envDefault:"C:\\Users\\Administrador\\golang\\ms-chat\\auth-service"`
}

var Envs = initConfig()

func initConfig() Config {
	if err := godotenv.Load(); err != nil {
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
