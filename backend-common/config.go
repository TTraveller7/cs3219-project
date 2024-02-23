package common

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var allowOrigins []string = []string{
	"http://user-service:8000",
	"http://auth-service:13704",
	"http://question-service:17001",
	"http://matching-service:7001",
	"http://reader-routine:11200",
	"http://socketio-matching-service:5200",
	"http://socketio-chat-service:5300",
	"http://socketio-collab-service:5400",
	"http://reverse-proxy:80",
}

func GetUrl(envKey string, defaultAddr string) string {
	if envUrl := os.Getenv(envKey); len(envUrl) != 0 {
		return envUrl
	}
	return defaultAddr
}

func Cors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = allowOrigins
	config.AllowCredentials = true

	return cors.New(config)
}
