package main

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var cache *redis.Client

func saveToCache(user *user, jwtString string, exp time.Time) bool {
	err := cache.Set(context.Background(), jwtString, user.Username, time.Now().Sub(exp)).Err()

	if err != nil {
		log.Error("Fail to save JWT of "+user.Username+":", err)
		return false
	}

	return true
}

func isInCache(username string, jwtString string) bool {
	val, err := cache.Get(context.Background(), jwtString).Result()

	if err != nil {
		log.Error("Fail to check if jwtString is in cache:", err)
		return false
	}

	if val == username {
		return true
	} else {
		return false
	}
}
