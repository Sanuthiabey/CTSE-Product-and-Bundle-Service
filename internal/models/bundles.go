package models

type Bundle struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Mood     string          `json:"mood"`
	Products []BundleProduct `json:"products,omitempty"`
}

type BundleProduct struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CreateBundleRequest struct {
	ID       string          `json:"id" binding:"required"`
	Name     string          `json:"name" binding:"required"`
	Mood     string          `json:"mood"`
	Products []BundleProduct `json:"products" binding:"required,dive,required"`
}
