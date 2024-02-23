package common

import (
	"fmt"
	"net/http"
	"time"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	SECRET = "sdfadfafawejonkfdsnvoqweorqrdaf"
)

type Middleware struct {
	cache  *redis.Client
	pool   *pgxpool.Pool
	log    *Logger
	SECRET string
}

func CreateMiddleware(cache *redis.Client, pool *pgxpool.Pool, log *Logger) *Middleware {
	return &Middleware{cache, pool, log, SECRET}
}

func (m *Middleware) RequireAuth(c *gin.Context) {
	// Get the cookie off request
	tokenString, err := c.Cookie("Authorization")

	abortWithErrorMsg := func() {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

	if err != nil {
		m.log.Error("", err)
		abortWithErrorMsg()
		return
	}

	// Decode / Validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(SECRET), nil
	})

	if err != nil {
		m.log.Error("", err)
		abortWithErrorMsg()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// If the Jwt expires, abort
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			m.log.Message("User unauthorized: JWT expired")
			abortWithErrorMsg()
			return
		}

		// If Jwt exists in blacklist, abort
		isBlocked := m.isInCache(claims["sub"].(string), tokenString)
		if isBlocked {
			m.log.Message("User unauthorized: invalid jwt")
			abortWithErrorMsg()
			return
		}

		// Find the user with token sub
		var username string = m.getUser(claims["sub"].(string))

		if username == "" {
			m.log.Message("User unauthorized: user with name " + claims["sub"].(string) + " not found")
			abortWithErrorMsg()
		}

		// Store username and expiration time into Context
		c.Set("Username", username)
		c.Set("Exp", int64(claims["exp"].(float64)))

		// Continue
		c.Next()

	} else {
		m.log.Message("User unauthorized: invalid token")
		abortWithErrorMsg()
	}
}

func (m *Middleware) isInCache(username string, jwtString string) bool {
	val, err := m.cache.Get(context.Background(), jwtString).Result()

	if err != nil {
		return false
	}

	if val == username {
		return true
	} else {
		return false
	}
}

func (m *Middleware) getUser(userName string) string {
	var name string = ""
	r := m.pool.QueryRow(context.Background(), "SELECT username FROM users WHERE username = $1::text", userName)

	err := r.Scan(&name)
	if err != nil {
		m.log.Error("Fail to find user with name "+userName, err)
		return ""
	}

	return name
}
