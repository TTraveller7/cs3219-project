package main

import (
	"fmt"

	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func createTables() {
	db.AutoMigrate(&Question{}, &Answer{})
}

func findQuestionById(questionId uint) (*Question, bool) {
	question := &Question{}
	if err := db.First(question, questionId).Error; err != nil {
		log.Error("", err)
		return question, false
	}

	return question, true
}

func questionIdList(difficulty string) []uint {
	// Retrieve size of questions list of the difficulty
	var size int64
	if err := db.Model(&Question{}).Where("difficulty = ?", difficulty).Count(&size).Error; err != nil {
		log.Error("Fail to retrieve size of question list", err)
		return make([]uint, 0)
	}

	list := make([]uint, size)
	for i := range list {
		list[i] = uint(i) + 1
	}

	shuffle(list)

	log.Message(fmt.Sprintf("List created: %#v", list))
	return list
}

func findAnswerByMatchIdAndQuestionId(matchId string, questionId uint) (*Answer, bool) {
	answer := &Answer{}

	// Take the last entry
	if err := db.Preload("Question").Where("match_id = ? AND question_id = ?", matchId, questionId).Last(answer).Error; err != nil {
		log.Error("", err)
		return answer, false
	}
	return answer, true
}

func saveAnswer(a *Answer) bool {
	if err := db.Create(a).Error; err != nil {
		log.Error("", err)
		return false
	}
	return true
}
