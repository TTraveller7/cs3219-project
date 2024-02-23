package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	SECRET = "sdfadfafawejonkfdsnvoqweorqrdaf"
)

var pool *pgxpool.Pool
var cache *redis.Client

func tryParseUserJwt(tokenString string) (string, bool) {
	claims, ok := decodeJwt(tokenString)
	if !ok {
		return "", false
	}

	// If the Jwt expires, abort
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		log.Message("User unauthorized: JWT expired")
		return "", false
	}

	// If Jwt exists in blacklist, abort
	username, _ := claims["sub"].(string)
	if isBlocked := isJwtInCache(username, tokenString); isBlocked {
		log.Message("User unauthorized: invalid jwt")
		return "", false
	}

	return username, true
}

func tryParseMatchId(tokenString string) (string, bool) {
	claims, ok := decodeJwt(tokenString)
	if !ok {
		return "", false
	}

	// If Jwt exists in blacklist, abort
	matchId, _ := claims["sub"].(string)
	if isBlocked := isJwtInCache(matchId, tokenString); isBlocked {
		log.Message("User unauthorized: invalid jwt")
		return "", false
	}

	// TODO: check if match is ended

	return matchId, true
}

func decodeJwt(tokenString string) (map[string]interface{}, bool) {
	// Decode / Validate it
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		log.Error("", err)
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || !ok {
		return nil, false
	}

	return claims, true
}

func isJwtInCache(username string, jwtString string) bool {
	val, err := cache.Get(mainContext, jwtString).Result()

	if err != nil {
		return false
	}

	if val == username {
		return true
	} else {
		return false
	}
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(SECRET), nil
}
