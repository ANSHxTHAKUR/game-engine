package model

import "time"

type UserResponse struct {
	UserID    string `json:"user_id"`
	Answer    string `json:"answer"`
	IsCorrect bool   `json:"is_correct"`
	Timestamp time.Time
}

type GameResult struct {
	WinnerID       string
	TotalCorrect   int64
	TotalIncorrect int64
	TimeTaken      time.Duration
}
