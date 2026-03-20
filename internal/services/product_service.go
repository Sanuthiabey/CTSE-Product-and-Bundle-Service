package services

import (
	"database/sql"
	"fmt"

	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/db"
	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/models"
)

func CreateProduct(product models.Product) error {

	_, err := db.DB.Exec(`
	INSERT INTO products
	(id, name, description, price, mood, category, image, rating, featured, stock)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.Mood,
		product.Category,
		product.Image,
		product.Rating,
		product.Featured,
		product.Stock,
	)

	return err
}

func GetAllProducts(limit int, offset int) ([]models.Product, error) {

	rows, err := db.DB.Query(`
	SELECT id,name,description,price,mood,category,image,rating,featured,stock
	FROM products
	ORDER BY id
	LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []models.Product

	for rows.Next() {

		var p models.Product

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Mood,
			&p.Category,
			&p.Image,
			&p.Rating,
			&p.Featured,
			&p.Stock,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func GetProductByID(id string) (*models.Product, error) {

	var p models.Product

	err := db.DB.QueryRow(`
	SELECT id,name,description,price,mood,category,image,rating,featured,stock
	FROM products WHERE id=$1
	`, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Mood,
		&p.Category,
		&p.Image,
		&p.Rating,
		&p.Featured,
		&p.Stock,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func UpdateProduct(id string, updated models.Product) error {

	result, err := db.DB.Exec(`
	UPDATE products SET
	name=$1,
	description=$2,
	price=$3,
	mood=$4,
	category=$5,
	image=$6,
	rating=$7,
	featured=$8,
	stock=$9
	WHERE id=$10
	`,
		updated.Name,
		updated.Description,
		updated.Price,
		updated.Mood,
		updated.Category,
		updated.Image,
		updated.Rating,
		updated.Featured,
		updated.Stock,
		id,
	)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func DeleteProduct(id string) error {

	result, err := db.DB.Exec(
		"DELETE FROM products WHERE id=$1",
		id,
	)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
