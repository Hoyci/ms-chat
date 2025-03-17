package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/hoyci/ms-chat/message-service/config"
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

func IsUserOnline(userID int) bool {
	count, err := GetClient().Get(context.Background(), fmt.Sprintf("connections:%d", userID)).Int()
	if err != nil && err != redis.Nil {
		log.Printf("An unexpected error occurred while checking user status: %v", err)
		return false
	}

	return count > 0
}
