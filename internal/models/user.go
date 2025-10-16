package models

import "time"

// UserCredentials креды пользователя
// @Description Запрос на регистрацию и аутентификацию
type UserCredentials struct {
	Login    string `json:"login" example:"Djon"`
	Password string `json:"password" example:"123456"`
}

// UserBalance баланс пользователя
// @Description Запрос на получения баланса
type UserBalance struct {
	Current   float64 `json:"current" example:"1000"`
	Withdrawn float64 `json:"withdrawn" example:"99"`
}

// UserWithdrawals вывод средств
// @Description Запрос на историю выводов средств
type UserWithdrawals struct {
	Order       string    `json:"order" example:"2377225624"`
	Sum         float64   `json:"sum" example:"10"`
	ProcessedAt time.Time `json:"processed_at" example:"2025-10-09T16:09:57+03:00"`
}
