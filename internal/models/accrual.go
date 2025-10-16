package models

import "github.com/google/uuid"

// AccrualRequest данные о заказе
// @Description Запрос в систему лояльности, для расчета баллов
type AccrualRequest struct {
	UserID uuid.UUID
	Order  string  `json:"order" example:"2377225624"`
	Goods  []Goods `json:"goods"`
}

// AccrualResponse данные о заказе
// @Description Ответ из системы лояльности, после расчета баллов
type AccrualResponse struct {
	Order   string  `json:"order" example:"2377225624"`
	Status  string  `json:"status" example:"PROCESSED"`
	Accrual float64 `json:"accrual" example:"100"`
}
