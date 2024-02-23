package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {
	// Get the cookie off request
	tokenString, err := c.Cookie("Authorization")

	abortWithErrorMsg := func() {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

	if err != nil {
		log.Error("", err)
		abortWithErrorMsg()
		return
	}

	// Decode / Validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(middleware.SECRET), nil
	})

	if err != nil {
		log.Error("", err)
		abortWithErrorMsg()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// If the Jwt expires, abort
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			log.Message("User unauthorized: JWT expired")
			abortWithErrorMsg()
			return
		}

		// If Jwt exists in blacklist, abort
		isBlocked := isInCache(claims["sub"].(string), tokenString)
		if isBlocked {
			log.Message("User unauthorized: invalid jwt")
			abortWithErrorMsg()
			return
		}

		// Find the user with token sub
		var u *user = getUser(claims["sub"].(string))

		if u == nil {
			log.Message("User unauthorized: user with name " + claims["sub"].(string) + " not found")
			abortWithErrorMsg()
		}

		// Store username and expiration time into Context
		c.Set("Username", u.Username)
		c.Set("Exp", int64(claims["exp"].(float64)))

		// Continue
		c.Next()

	} else {
		log.Message("User unauthorized: invalid token")
		abortWithErrorMsg()
	}
}
