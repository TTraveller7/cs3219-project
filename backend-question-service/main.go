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

	questionRoutes := engine.Group("/question")
	{
		questionRoutes.GET("/:questionId", getQuestionWithIdHandler)
		questionRoutes.GET("/curr", common.ExtractMatchId, getCurrentQuestionHandler)
		questionRoutes.GET("/next", common.ExtractMatchId, getNextQuestionHandler)
		questionRoutes.GET("/start", common.ExtractMatchId, getStartingQuestionHandler)
	}

	answerRoutes := engine.Group("/answer")
	{
		answerRoutes.GET("", getAnswerHandler)
		answerRoutes.POST("/create", common.ExtractMatchId, saveAnswerHandler)
	}

	return engine
}

func main() {
	log = common.CreateLogger("question service")
	mainContext = context.Background()

	log.Message("Connecting to postgres and cache")

	// Setup db
	initDb()

	// Setup cache
	questionCache = common.CreateCache(getQuestionCacheUrl(), log)
	defer common.CloseCache(questionCache)

	matchCache = common.CreateCache(getMatchRedisUrl(), log)
	defer common.CloseCache(matchCache)

	// Setup routes
	engine := setupEngine()

	// Start server
	log.Message("Starting server...")
	if err := engine.Run(ADDRESS); err != nil {
		log.Error("Fail to start server:", err)
	}
}
