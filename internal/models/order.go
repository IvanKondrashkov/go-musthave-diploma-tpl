package models

import "time"

// Order заказ
// @Description Запрос на создание заказа
type Order struct {
	Number     string    `json:"number" example:"2377225624"`
	Status     string    `json:"status" example:"NEW"`
	Accrual    float64   `json:"accrual" example:"100"`
	UploadedAt time.Time `json:"uploaded_at" example:"2025-10-09T16:09:57+03:00"`
}
