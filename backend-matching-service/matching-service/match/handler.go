package match

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/backend-common/common"
	"github.com/gin-gonic/gin"
)

func abortWithErrorMsg(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
}

// Handle match request

func BeginMatching(c *gin.Context) {
	abortWithErrorMsg := func(msg string) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
	}
	ctx := context.Background()

	// Extract MatchRequest from request body
	matchRequest := &MatchRequest{}
	if err := c.ShouldBind(matchRequest); err != nil && common.DifficultySet()[matchRequest.Difficulty] {
		log.Error("Fail to bind match request", err)
		abortWithErrorMsg("Fail to start matching")
		return
	}

	username := matchRequest.Username

	// Check with redis: if username is already waiting / already has a match
	if !isUserInCache(ctx, username) {
		saveUser(ctx, &User{
			Username: username,
			Status:   common.USER_STATUS_IDLE,
		})
	} else {
		u, err := getUser(ctx, username)
		if err != nil {
			log.Message("Fail to get user with username " + username)
			abortWithErrorMsg("Fail to bind match request")
			return
		}

		// If status of the user is not idle, the user cannot create match
		if u.Status != common.USER_STATUS_IDLE {
			log.Message("User with status" + u.Status + " cannot create match")
			abortWithErrorMsg("User is in a matching queue or has an ongoing match.")
			return
		}
	}

	// Assign uuid and uuid to matchRequest
	matchRequest.MatchRequestId = getNewUuid()
	matchRequest.CreatedAt = time.Now()

	// update user status
	statusStr := getPendingStatus(matchRequest.Difficulty)
	if err := setUser(ctx, username, "Status", statusStr,
		"LastMatchRequestId", matchRequest.MatchRequestId); err != nil {
		log.Error("Fail to set user status", err)
		abortWithErrorMsg("Fail to start matching")
		return
	}

	// Serialize match request
	matchRequestJson, err := json.Marshal(matchRequest)
	if err != nil {
		log.Error("Fail to serialize match request to json", err)
		abortWithErrorMsg("Fail to start matching")
		return
	}

	// write to kafka
	var isWriteSuccussful bool
	switch matchRequest.Difficulty {
	case common.DIFFICULTY_EASY:
		isWriteSuccussful = WriteToEasyQueue(matchRequestJson)
	case common.DIFFICULTY_MEDIUM:
		isWriteSuccussful = WriteToMediumQueue(matchRequestJson)
	case common.DIFFICULTY_HARD:
		isWriteSuccussful = WriteToHardQueue(matchRequestJson)
	}
	if !isWriteSuccussful {
		log.Message("Fail to create match: fail to write to kafka")
		abortWithErrorMsg("Fail to create match")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Match request sent"})
}

// Handle user

func ExtractUserWithName(c *gin.Context) {
	handleRequestCtx, cancelHandleRequest := context.WithCancel(mainContext)
	defer cancelHandleRequest()

	username := c.Param("username")
	user, err := getUser(handleRequestCtx, username)
	if err != nil {
		log.Error("Fail to extract user with name "+username, err)
		abortWithErrorMsg(c, "Fail to find user with name "+username)
		return
	}

	c.Set("user", user)
}

func GetUserWithName(c *gin.Context) {
	user, _ := c.Keys["user"].(*User)

	handleMatchIdCookie(c)

	c.JSON(http.StatusOK, gin.H{"user": *user, "message": "Successfully retrieved user"})
}

func ChangeUserStatusToIdle(c *gin.Context) {
	handleRequestCtx, cancelHandleRequest := context.WithCancel(mainContext)
	defer cancelHandleRequest()

	user, _ := c.Keys["user"].(*User)

	if err := updateUserStatus(handleRequestCtx, user.Username, "", common.USER_STATUS_IDLE); err != nil {
		log.Error("Fail to change user status to idle:", err)
		abortWithErrorMsg(c, "Fail to change user status")
		return
	}

	user.Status = common.USER_STATUS_IDLE
	c.Set("user", user)
	handleMatchIdCookie(c)

	c.JSON(http.StatusOK, gin.H{"message": "User status set to idle"})
}

// Handle match id cookie

func handleMatchIdCookie(c *gin.Context) {
	user, _ := c.Keys["user"].(*User)

	c.SetSameSite(http.SameSiteLaxMode)

	if user.Status == common.USER_STATUS_MATCHED || user.Status == common.USER_STATUS_LEAVING {
		c.SetCookie("MatchId", createMatchIdCookie(user.MatchId), 3600*24*30, "/", "", false, true)
	} else {
		c.SetCookie("MatchId", "", -1, "/", "", false, true)
	}
}

// Handle match

func ExtractMatch(c *gin.Context) {
	handleRequestCtx, cancelHandleRequest := context.WithCancel(mainContext)
	defer cancelHandleRequest()

	matchId := c.Param("matchId")
	match, err := getMatch(handleRequestCtx, matchId)
	if err != nil {
		log.Error("Fail to extract match with id "+matchId, err)
		abortWithErrorMsg(c, "Fail to find match with id "+matchId)
		return
	}

	c.Set("match", match)
}

func GetMatchWithId(c *gin.Context) {
	match, ok := c.Keys["match"].(*Match)
	if !ok {
		log.Error("Fail to retrieve match:", errors.New("fail to get match from context"))
		abortWithErrorMsg(c, "Fail to retrieve match")
		return
	}

	c.JSON(http.StatusOK, gin.H{"match": *match, "message": "Successfully retrieved match"})
}

func EndMatchWithId(c *gin.Context) {
	handleRequestCtx, cancelHandleRequest := context.WithCancel(mainContext)
	defer cancelHandleRequest()

	match, ok := c.Keys["match"].(*Match)
	if !ok {
		log.Error("Fail to retrieve match:", errors.New("fail to get match from context"))
		abortWithErrorMsg(c, "Fail to retrieve match")
		return
	}

	if match.IsEnded {
		c.JSON(http.StatusOK, gin.H{"message": "Match already ended"})
		return
	}

	match.IsEnded = true
	match.EndedAt = time.Now()
	saveMatch(handleRequestCtx, match)

	c.JSON(http.StatusOK, gin.H{"message": "Match successfully ended"})
}
