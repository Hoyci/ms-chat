package redis

import (
	"context"
	"log"

	"github.com/hoyci/ms-chat/ws-service/config"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func Init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.Envs.RedisAddr,
		Password: config.Envs.RedisPassword,
		DB:       config.Envs.RedisDB,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

func GetClient() *redis.Client {
	return redisClient
}
