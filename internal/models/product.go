package models

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Mood        string  `json:"mood"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
	Rating      float64 `json:"rating"`
	Featured    bool    `json:"featured"`
	Stock       int     `json:"stock"`
}
