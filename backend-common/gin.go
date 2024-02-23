package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Extract match id from header. Abort if match id is empty.
func ExtractMatchId(c *gin.Context) {
	matchId := c.GetHeader(MATCH_ID_HEADER_KEY)
	if len(matchId) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseBodyWithError("Fail to extract match id"))
		return
	}
	c.Set(MATCH_ID_CONTEXT_KEY, matchId)
}

func ExtractUsername(c *gin.Context) {
	c.Set(USERNAME_CONTEXT_KEY, c.GetHeader(USERNAME_HEADER_KEY))
}

func Healthcheck(c *gin.Context) {
	c.Status(http.StatusAccepted)
}

func EngineWithHealthcheck() *gin.Engine {
	engine := gin.New()
	engine.Use(
		gin.LoggerWithConfig(
			gin.LoggerConfig{
				SkipPaths: []string{HEALTHCHECK_URL},
			},
		),
		gin.Recovery(),
	)
	engine.GET(HEALTHCHECK_URL, Healthcheck)

	return engine
}
