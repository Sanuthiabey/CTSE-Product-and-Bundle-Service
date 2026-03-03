package models

type Bundle struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Mood string `json:"mood"`
}

type BundleProduct struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
