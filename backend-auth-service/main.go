package main

import (
	"context"

	"github.com/backend-common/common"
	"github.com/gin-gonic/gin"
)

var (
	log         *common.Logger
	mainContext context.Context
)

func setupEngine() *gin.Engine {
	engine := common.EngineWithHealthcheck()

	engine.GET("/auth", authorize)

	return engine
}

func main() {
	log = common.CreateLogger("auth service")
	mainContext = context.Background()

	log.Message("Connecting to postgres and redis")

	// Set up postgres database connection pool
	dbpool := common.CreateDbPool(getPostgresUrl(), log)
	defer common.CloseDbpool(dbpool)
	pool = dbpool

	// Set up redis
	redis_client := common.CreateCache(getRedisUrl(), log)
	defer common.CloseCache(redis_client)
	cache = redis_client

	engine := setupEngine()

	log.Message("Starting server...")

	err := engine.Run(ADDRESS)
	if err != nil {
		log.Error("Fail to start server:", err)
	}
}
