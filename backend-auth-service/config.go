package main

import (
	"github.com/backend-common/common"
)

const (
	DEFAULT_URL = "postgres://postgres:123@localhost:5525/dev"
	ADDRESS     = ":13704"

	REDIS_ADDR = "localhost:6379"
)

func getPostgresUrl() string {
	return common.GetUrl(common.ENV_POSTGRES_URL, DEFAULT_URL)
}

func getRedisUrl() string {
	return common.GetUrl(common.ENV_BLACKLIST_REDIS_URL, REDIS_ADDR)
}
