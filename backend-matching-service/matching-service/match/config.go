package match

import (
	"github.com/backend-common/common"
)

const (
	DEFAULT_URL = "postgres://postgres:123@localhost:5000/dev"
	ADDRESS     = ":7001"

	KAFKA_URL         = "localhost:9092"
	KAFKA_RETRY_TIMES = 3

	USER_REDIS_ADDR  = "localhost:6380"
	MATCH_REDIS_ADDR = "localhost:6381"
)

func getPostgresUrl() string {
	return common.GetUrl(common.ENV_POSTGRES_URL, DEFAULT_URL)
}

func getKafkaUrl() string {
	return common.GetUrl(common.ENV_KAFKA_URL, KAFKA_URL)
}

func getUserRedisUrl() string {
	return common.GetUrl(common.ENV_USER_REDIS_URL, USER_REDIS_ADDR)
}

func getMatchRedisUrl() string {
	return common.GetUrl(common.ENV_MATCH_REDIS_URL, MATCH_REDIS_ADDR)
}
