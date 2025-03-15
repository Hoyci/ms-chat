package config

import (
	"log"
	"os"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Port          int    `env:"PORT" envDefault:"8082"`
	Environment   string `env:"ENVIRONMENT" envDefault:"development"`
	DatabaseURL   string `env:"DATABASE_URL" envDefault:"mongodb://root:example@localhost:27017/"`
	DatabaseName  string `env:"DATABASE_NAME" envDefault:"admin"`
	RabbitMQURL   string `env:"RABBITMQ_URL" envDefault:"amqp://user:password@localhost:5672/"`
	RedisAddr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:"password"`
	RedisDB       int    `env:"REDIS_DB" envDefault:"0"`
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
