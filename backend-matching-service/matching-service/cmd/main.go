package main

import (
	"context"
	"matching-service/match"

	"github.com/backend-common/common"
	"github.com/gin-gonic/gin"
)

var (
	log         *common.Logger
	mainContext context.Context
)

func setupEngine() *gin.Engine {
	engine := common.EngineWithHealthcheck()

	// handle match request
	engine.POST("/match/create", match.BeginMatching)

	// handle user
	userRoutes := engine.Group("/user/:username")
	{
		userRoutes.Use(match.ExtractUserWithName)
		userRoutes.GET("", match.GetUserWithName)
		userRoutes.PUT("/toidle", match.ChangeUserStatusToIdle)
	}

	// handle match
	matchRoutes := engine.Group("/match/:matchId")
	{
		matchRoutes.Use(match.ExtractMatch)
		matchRoutes.GET("", match.GetMatchWithId)
		matchRoutes.PUT("/end", match.EndMatchWithId)
	}

	return engine
}

func main() {
	match.TryConnectToKafka()

	log = common.CreateLogger("matching service")
	mainContext = context.Background()
	match.Setup(mainContext, log)

	// set up kafka
	match.SetupKafka()
	defer match.CloseKafka()

	// Set up cache and db
	log.Message("Connecting to postgres, kafka, and redis")

	match.ConnectCacheAndDb()
	defer match.CloseCacheAndDb()

	engine := setupEngine()

	log.Message("Starting server...")

	err := engine.Run(match.ADDRESS)
	if err != nil {
		log.Error("Fail to start server:", err)
	}
}
