package models

// Goods товары
type Goods struct {
	Description string  `json:"description" example:"Кофе машина"`
	Price       float64 `json:"price" example:"1000"`
}
