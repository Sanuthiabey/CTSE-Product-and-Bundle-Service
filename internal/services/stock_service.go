package services

import (
	"database/sql"
	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/db"
	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/models"
)

// VALIDATE STOCK
func ValidateStock(productID string, qty int) (map[string]interface{}, error) {

	var stock int

	err := db.DB.QueryRow(
		"SELECT stock FROM products WHERE id=$1",
		productID,
	).Scan(&stock)

	if err == sql.ErrNoRows {
		return map[string]interface{}{
			"available": false,
			"message":   "Product not found",
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"available": stock >= qty,
		"stock":     stock,
	}, nil
}

// SET STOCK (ADMIN)
func SetStock(productID string, stock int) error {

	_, err := db.DB.Exec(
		"UPDATE products SET stock = $1 WHERE id = $2",
		stock,
		productID,
	)

	return err
}

// REDUCE STOCK
func ReduceStock(items []models.StockItem) error {

	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	for _, item := range items {

		var stock int

		err := tx.QueryRow(
			"SELECT stock FROM products WHERE id=$1",
			item.ProductID,
		).Scan(&stock)

		if err != nil || stock < item.Quantity {
			tx.Rollback()
			return err
		}
	}

	for _, item := range items {
		_, err := tx.Exec(
			"UPDATE products SET stock = stock - $1 WHERE id=$2",
			item.Quantity,
			item.ProductID,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
