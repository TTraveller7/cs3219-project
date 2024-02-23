package main

import (
	"github.com/backend-common/common"
)

const (
	POSTGRES_ADDR = "postgres://postgres:123@localhost:5000/dev"
	ADDRESS       = ":8000"
	REDIS_ADDR    = "localhost:6379"
)

func getPostgresUrl() string {
	return common.GetUrl(common.ENV_POSTGRES_URL, POSTGRES_ADDR)
}

func getRedisUrl() string {
	return common.GetUrl(common.ENV_BLACKLIST_REDIS_URL, REDIS_ADDR)
}
