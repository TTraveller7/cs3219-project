package main

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	Difficulty  string `json:"difficulty"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Answer struct {
	gorm.Model
	MatchID    string   `json:"matchId"`
	Code       string   `json:"code"`
	QuestionID uint     `json:"questionId" binding:"required"`
	Question   Question `json:"question"`
}
