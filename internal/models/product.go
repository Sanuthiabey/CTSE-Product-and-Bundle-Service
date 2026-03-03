package models

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Mood  string  `json:"mood"`
	Stock int     `json:"stock"`
}
