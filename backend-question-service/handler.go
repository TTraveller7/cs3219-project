package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/backend-common/common"
	"github.com/gin-gonic/gin"
)

func getStartingQuestionHandler(c *gin.Context) {
	// Init handler context
	handlerContext, cancelHandle := context.WithCancel(mainContext)
	defer cancelHandle()

	abortWithErrMsg := func() {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ResponseBodyWithError("Fail to get next question"))
	}

	// Get match id from context
	matchId := c.GetString(common.MATCH_ID_CONTEXT_KEY)

	numOfRetries := 10
	if ok := tryLock(handlerContext, matchId, numOfRetries); !ok {
		log.Error("", fmt.Errorf("fails to fetch lock in %d retries", numOfRetries))
		abortWithErrMsg()
		return
	}
	defer unlock(handlerContext, matchId)

	// check if qid list exists and is not empty
	isQidListEmpty := isQuestionIdListEmpty(handlerContext, matchId)

	// If list is empty, create new question id list based on difficulty
	if isQidListEmpty {
		if ok, err := initQidListAndPointer(handlerContext, matchId); !ok || err != nil {
			log.Error("", errors.New("fail to init question id and pointer for matchId "+matchId))
			abortWithErrMsg()
			return
		}

		// Update pointer to 2
		savePointer(handlerContext, matchId, INITIAL_POINTER+1)
	}

	qid, err := getQidByPointer(handlerContext, matchId, INITIAL_POINTER)
	if err != nil {
		log.Error("", err)
		abortWithErrMsg()
		return
	}

	// Get next question object
	question, _ := findQuestionById(qid)

	// Generate success response
	c.JSON(http.StatusOK, common.ResponseBodyWithMessage("Successfully retrieved question", "question", question))
}

func getNextQuestionHandler(c *gin.Context) {
	// Init handler context
	handlerContext, cancelHandle := context.WithCancel(mainContext)
	defer cancelHandle()

	abortWithErrMsg := func() {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ResponseBodyWithError("Fail to fetch the next question. Please try again later."))
	}

	// Get match id from context
	matchId := c.GetString(common.MATCH_ID_CONTEXT_KEY)

	if ok := lock(handlerContext, matchId); !ok {
		log.Error("", fmt.Errorf("processing another /next/question request. aborting"))
		abortWithErrMsg()
		return
	}
	defer unlock(handlerContext, matchId)

	// check if qid list exists and is not empty
	isQidListEmpty := isQuestionIdListEmpty(handlerContext, matchId)

	// If list is empty, abort with BadRequest
	if isQidListEmpty {
		log.Message("The question list is empty.")
		abortWithErrMsg()
		return
	}

	ptr, err := getPointer(handlerContext, matchId)
	if err != nil {
		log.Error("Pointer with match id "+matchId+" not found", err)
		abortWithErrMsg()
		return
	}

	qid, err := getQidByPointer(handlerContext, matchId, ptr)
	if err != nil {
		log.Error("", err)
		abortWithErrMsg()
		return
	}

	// Get next question object
	question, _ := findQuestionById(qid)

	// Update pointer
	savePointer(handlerContext, matchId, ptr+1)

	// Generate success response
	c.JSON(http.StatusOK, common.ResponseBodyWithMessage("Successfully retrieved question", "question", question))
}

func getCurrentQuestionHandler(c *gin.Context) {
	// Init handler context
	handlerContext, cancelHandle := context.WithCancel(mainContext)
	defer cancelHandle()

	abortWithErrMsg := func() {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ResponseBodyWithError("Fail to fetch the current question. Please try again later."))
	}

	// Get match id from context
	matchId := c.GetString(common.MATCH_ID_CONTEXT_KEY)

	ptr, err := getPointer(handlerContext, matchId)
	if err != nil {
		log.Error("Pointer with match id "+matchId+" not found", err)
		abortWithErrMsg()
		return
	}

	qid, err := getQidByPointer(handlerContext, matchId, ptr-1)
	if err != nil {
		log.Error("", err)
		abortWithErrMsg()
		return
	}

	// Get question object
	question, _ := findQuestionById(qid)

	// Generate success response
	c.JSON(http.StatusOK, common.ResponseBodyWithMessage("Successfully retrieved question", "question", question))
}

func getQuestionWithIdHandler(c *gin.Context) {
	// Extract question id from url param
	questionId, _ := strconv.ParseUint(c.Param("questionId"), 10, 32)

	// Try find question in db
	if question, exist := findQuestionById(uint(questionId)); !exist {
		c.JSON(http.StatusBadRequest, common.ResponseBodyWithError("Fail to retrieve question"))
	} else {
		c.JSON(http.StatusOK, common.ResponseBodyWithMessage("Successfully retrieved question", "question", question))
	}
}

// Answer routes

func getAnswerHandler(c *gin.Context) {
	abortWithErrMsg := func() {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ResponseBodyWithError("Fail to get answer"))
	}

	// Extract matchId and question id from query param
	matchId := c.Query(common.MATCH_ID_PARAM_KEY)
	questionIdStr := c.Query(common.QUESTION_ID_PARAM_KEY)
	questionId, err := common.StringToUint(questionIdStr)
	if len(matchId) == 0 || err != nil {
		log.Error("", errors.New("matchId is empty or question id is invalid"))
		abortWithErrMsg()
		return
	}

	if answer, success := findAnswerByMatchIdAndQuestionId(matchId, questionId); !success {
		log.Message("Fail to retrieve answer")
		abortWithErrMsg()
	} else {
		c.JSON(http.StatusOK, common.ResponseBodyWithMessage("Successfully retrieved answer", "answer", *answer))
	}
}

func saveAnswerHandler(c *gin.Context) {
	abortWithErrMsg := func() {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ResponseBodyWithError("Fail to save answer"))
	}

	matchId := c.GetString(common.MATCH_ID_CONTEXT_KEY)

	// Bind answer from request body
	answer := &Answer{}
	if err := c.ShouldBind(answer); err != nil {
		log.Error("", err)
		abortWithErrMsg()
		return
	}

	answer.MatchID = matchId
	if success := saveAnswer(answer); !success {
		log.Message("Fail to save answer")
		abortWithErrMsg()
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Successfully saved answer"})
	}
}
