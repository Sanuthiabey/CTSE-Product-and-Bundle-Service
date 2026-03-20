package services

import (
	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/db"
	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/models"
)

// CREATE BUNDLE
func CreateBundle(req models.CreateBundleRequest) (models.Bundle, error) {

	tx, err := db.DB.Begin()
	if err != nil {
		return models.Bundle{}, err
	}

	_, err = tx.Exec(`
		INSERT INTO bundles (id, name, mood)
		VALUES ($1,$2,$3)
	`, req.ID, req.Name, req.Mood)

	if err != nil {
		tx.Rollback()
		return models.Bundle{}, err
	}

	for _, p := range req.Products {
		_, err := tx.Exec(`
			INSERT INTO bundle_products (bundle_id, product_id, quantity)
			VALUES ($1,$2,$3)
		`, req.ID, p.ProductID, p.Quantity)

		if err != nil {
			tx.Rollback()
			return models.Bundle{}, err
		}
	}

	tx.Commit()

	return models.Bundle{
		ID:       req.ID,
		Name:     req.Name,
		Mood:     req.Mood,
		Products: req.Products,
	}, nil
}

// GET BUNDLES
func GetBundles() ([]models.Bundle, error) {

	rows, err := db.DB.Query(`SELECT id,name,mood FROM bundles`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bundles []models.Bundle

	for rows.Next() {
		var b models.Bundle
		rows.Scan(&b.ID, &b.Name, &b.Mood)

		productRows, _ := db.DB.Query(`
			SELECT product_id, quantity FROM bundle_products WHERE bundle_id=$1
		`, b.ID)

		var products []models.BundleProduct
		for productRows.Next() {
			var p models.BundleProduct
			productRows.Scan(&p.ProductID, &p.Quantity)
			products = append(products, p)
		}

		b.Products = products
		bundles = append(bundles, b)
	}

	return bundles, nil
}
