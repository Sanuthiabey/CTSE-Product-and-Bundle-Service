package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/db"
	grpcServer "github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/grpc"
	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/models"
)

func main() {

	// Connect to Neon PostgreSQL
	db.Connect()

	// Start gRPC server in background
	go grpcServer.StartGRPCServer()

	r := gin.Default()

	// ==============================
	// HEALTH CHECK
	// ==============================
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "product-and-bundle-service",
			"status":  "running",
		})
	})

	// ==============================
	// CREATE PRODUCT
	// ==============================
	r.POST("/products", func(c *gin.Context) {
		var product models.Product

		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.DB.Exec(
			`INSERT INTO products 
			(id, name, description, price, mood, category, image, rating, featured, stock)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
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

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, product)
	})

	// ==============================
	// GET ALL PRODUCTS
	// ==============================
	r.GET("/products", func(c *gin.Context) {

		rows, err := db.DB.Query(`
			SELECT id, name, description, price, mood, category, image, rating, featured, stock
			FROM products
		`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
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
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			products = append(products, p)
		}

		c.JSON(http.StatusOK, products)
	})

	// ==============================
	// GET PRODUCT BY ID
	// ==============================
	r.GET("/products/:id", func(c *gin.Context) {

		id := c.Param("id")

		var p models.Product

		err := db.DB.QueryRow(`
			SELECT id, name, description, price, mood, category, image, rating, featured, stock
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
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, p)
	})

	// ==============================
	// UPDATE PRODUCT
	// ==============================
	r.PUT("/products/:id", func(c *gin.Context) {

		id := c.Param("id")

		var updated models.Product
		if err := c.ShouldBindJSON(&updated); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		updated.ID = id
		c.JSON(http.StatusOK, updated)
	})

	// ==============================
	// DELETE PRODUCT
	// ==============================
	r.DELETE("/products/:id", func(c *gin.Context) {

		id := c.Param("id")

		result, err := db.DB.Exec("DELETE FROM products WHERE id=$1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
	})

	// ==============================
	// START HTTP SERVER
	// ==============================
	r.Run(":8080")
}
