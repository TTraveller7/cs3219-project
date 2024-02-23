package common

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

func CreateCache(redisUrl string, log *Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "",
		DB:       0,
	})

	if err := client.Exists(context.Background(), "hello").Err(); err != nil {
		client.Close()
		log.Error("fail to connect to redis at "+redisUrl+":", err)
		os.Exit(1)
	}

	return client
}

func CloseCache(cache *redis.Client) error {
	return cache.Close()
}
