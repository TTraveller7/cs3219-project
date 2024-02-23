package main

import (
	common "github.com/backend-common/common"
	"github.com/gin-gonic/gin"
)

var log *common.Logger
var middleware *common.Middleware

func setupEngine() *gin.Engine {
	var engine *gin.Engine = common.EngineWithHealthcheck()

	// Routes
	engine.POST("/user/create", CreateUser)
	engine.POST("/login", AuthUser)
	engine.POST("/logout", middleware.RequireAuth, LogoutUser)
	engine.POST("/user/delete", middleware.RequireAuth, DeleteUser)
	engine.POST("/validate", middleware.RequireAuth, ValidateUser)
	engine.POST("/user/changepwd", middleware.RequireAuth, ChangePassword)

	return engine
}

func main() {
	// set up log
	log = common.CreateLogger("user-service")

	// Set up postgres database connection pool
	dbpool := common.CreateDbPool(getPostgresUrl(), log)
	defer common.CloseDbpool(dbpool)
	pool = dbpool

	// Set up redis
	redis_client := common.CreateCache(getRedisUrl(), log)
	defer common.CloseCache(redis_client)
	cache = redis_client

	// set up middleware
	middleware = common.CreateMiddleware(cache, pool, log)

	// Configure server
	var engine *gin.Engine = setupEngine()

	// Run server
	if err := engine.Run(ADDRESS); err != nil {
		log.Error("Fail to start up server:", err)
	}
}
