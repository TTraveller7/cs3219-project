package main

import (
	"context"

	"matching-service/match"

	"github.com/backend-common/common"
)

func main() {
	match.TryConnectToKafka()

	log := common.CreateLogger("reader routine")
	mainContext := context.Background()
	match.Setup(mainContext, log)

	quit := make(chan int, 1)

	// Set up cache and db
	log.Message("Connecting to postgres, kafka, and redis")
	match.ConnectCacheAndDb()
	defer match.CloseCacheAndDb()

	// Start consumers
	log.Message("Starting reader routines...")
	match.InitConsumers(quit)
	match.StartConsumers(mainContext)
	defer match.CloseConsumers()

	// Blocks
	<-quit
}
