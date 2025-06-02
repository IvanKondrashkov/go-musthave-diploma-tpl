package models

import "github.com/google/uuid"

type AccrualRequest struct {
	UserID uuid.UUID
	Order  string  `json:"order"`
	Goods  []Goods `json:"goods"`
}

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
