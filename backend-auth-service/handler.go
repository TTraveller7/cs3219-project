package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func authorize(c *gin.Context) {
	abortWithErrorMsg := func() {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

	// Extract username

	userCookieStr, err := c.Cookie("Authorization")
	if err != nil {
		log.Error("", err)
		abortWithErrorMsg()
		return
	}

	if username, isValid := tryParseUserJwt(userCookieStr); isValid {
		c.Header("X-Username", username)
	} else {
		abortWithErrorMsg()
		return
	}

	// Try extract matchId

	matchCookieStr, err := c.Cookie("MatchId")
	if err != nil {
		c.Status(http.StatusOK)
		return
	}

	if matchId, isValid := tryParseMatchId(matchCookieStr); isValid {
		c.Header("X-MatchId", matchId)
	} else {
		c.Header("X-MatchId", "")
	}

	c.Status(http.StatusOK)
}
