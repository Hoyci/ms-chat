package config

import (
	"log"

	"github.com/caarlos0/env"
)

type Config struct {
	Port                 int    `env:"PORT" envDefault:"8081"`
	Environment          string `env:"ENVIRONMENT" envDefault:"development"`
	RabbitMQURL          string `env:"RABBITMQ_URL" envDefault:"amqp://user:password@localhost:5672/"`
	RedisAddr            string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisPassword        string `env:"REDIS_PASSWORD" envDefault:"password"`
	RedisDB              int    `env:"REDIS_DB" envDefault:"0"`
	PersistenceQueueName string `env:"PERSISTENCE_QUEUE_NAME" envDefault:"persistence_queue"`
	BroadcastQueueName   string `env:"BROADCAST_QUEUE_NAME" envDefault:"broadcast_queue"`
	UserEventsQueueName  string `env:"USER_EVENTS_QUEUE_NAME" envDefault:"user_events_queue"`
}

var Envs = initConfig()

func initConfig() Config {
	var cfg Config
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	return cfg
}
