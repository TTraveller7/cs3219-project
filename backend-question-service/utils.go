package main

import (
	"math/rand"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func shuffle(list []uint) {
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(list), func(i int, j int) {
		temp := list[i]
		list[i] = list[j]
		list[j] = temp
	})
}

func QuestionListKey(matchId string) string {
	return "key-" + matchId
}

func PointerKey(matchId string) string {
	return "pointer-" + matchId
}

func LockKey(matchId string) string {
	return "lock-" + matchId
}

func initDb() {
	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN: getPostgresUrl(),
	}), &gorm.Config{})
	if err != nil {
		log.Error("", err)
		os.Exit(1)
	}
	db = database
	createTables()
	saveDummyData()
}
