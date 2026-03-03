package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
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
	// CORS CONFIG (VERY IMPORTANT)
	// ==============================
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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

		mood := c.Query("mood")
		category := c.Query("category")
		search := c.Query("search")
		featured := c.Query("featured")
		sort := c.Query("sort")

		query := `
        SELECT id, name, description, price, mood, category, image, rating, featured, stock
        FROM products
        WHERE 1=1
    `
		args := []interface{}{}
		argID := 1

		if mood != "" {
			query += " AND LOWER(mood) = LOWER($" + fmt.Sprint(argID) + ")"
			args = append(args, mood)
			argID++
		}

		if category != "" {
			query += " AND LOWER(category) = LOWER($" + fmt.Sprint(argID) + ")"
			args = append(args, category)
			argID++
		}

		if search != "" {
			query += " AND (LOWER(name) LIKE LOWER($" + fmt.Sprint(argID) + ") OR LOWER(description) LIKE LOWER($" + fmt.Sprint(argID) + "))"
			args = append(args, "%"+search+"%")
			argID++
		}

		if featured == "true" {
			query += " AND featured = true"
		}

		switch sort {
		case "price-asc":
			query += " ORDER BY price ASC"
		case "price-desc":
			query += " ORDER BY price DESC"
		case "rating":
			query += " ORDER BY rating DESC"
		default:
			query += " ORDER BY id"
		}

		rows, err := db.DB.Query(query, args...)
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
