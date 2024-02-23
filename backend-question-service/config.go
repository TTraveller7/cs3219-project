package main

import (
	"github.com/backend-common/common"
)

const (
	DEFAULT_URL = "postgres://postgres:123@localhost:5000/dev"
	ADDRESS     = ":17001"

	QUESTION_CACHE_ADDR = "localhost:6382"
	MATCH_REDIS_ADDR    = "localhost:6381"

	MATCH_DIFFICULTY_KEY = "Difficulty"

	INITIAL_POINTER = 1
)

func getPostgresUrl() string {
	return common.GetUrl(common.ENV_POSTGRES_URL, DEFAULT_URL)
}

func getQuestionCacheUrl() string {
	return common.GetUrl(common.ENV_QUESTION_REDIS_URL, QUESTION_CACHE_ADDR)
}

func getMatchRedisUrl() string {
	return common.GetUrl(common.ENV_MATCH_REDIS_URL, MATCH_REDIS_ADDR)
}
