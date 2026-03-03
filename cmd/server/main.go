package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/db"
	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/models"
)

func main() {

	db.Connect()

	r := gin.Default()

	// ------------------------
	// HEALTH CHECK
	// ------------------------
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "product-and-bundle-service",
			"status":  "running",
		})
	})

	// ==============================
	// PRODUCT ROUTES
	// ==============================

	r.POST("/products", func(c *gin.Context) {
		var product models.Product

		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.DB.Exec(
			"INSERT INTO products (id, name, price, mood, stock) VALUES ($1, $2, $3, $4, $5)",
			product.ID, product.Name, product.Price, product.Mood, product.Stock,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, product)
	})

	r.GET("/products", func(c *gin.Context) {
		rows, err := db.DB.Query("SELECT id, name, price, mood, stock FROM products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		products := []models.Product{}

		for rows.Next() {
			var p models.Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Mood, &p.Stock); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			products = append(products, p)
		}

		c.JSON(http.StatusOK, products)
	})

	r.GET("/products/:id", func(c *gin.Context) {
		id := c.Param("id")

		var p models.Product

		err := db.DB.QueryRow(
			"SELECT id, name, price, mood, stock FROM products WHERE id=$1",
			id,
		).Scan(&p.ID, &p.Name, &p.Price, &p.Mood, &p.Stock)

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

	r.PUT("/products/:id", func(c *gin.Context) {
		id := c.Param("id")

		var updated models.Product
		if err := c.ShouldBindJSON(&updated); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := db.DB.Exec(
			"UPDATE products SET name=$1, price=$2, mood=$3, stock=$4 WHERE id=$5",
			updated.Name, updated.Price, updated.Mood, updated.Stock, id,
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
	// BUNDLE ROUTES
	// ==============================

	r.POST("/bundles", func(c *gin.Context) {

		type CreateBundleRequest struct {
			ID       string                 `json:"id"`
			Name     string                 `json:"name"`
			Mood     string                 `json:"mood"`
			Products []models.BundleProduct `json:"products"`
		}

		var req CreateBundleRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Insert bundle
		_, err := db.DB.Exec(
			"INSERT INTO bundles (id, name, mood) VALUES ($1, $2, $3)",
			req.ID, req.Name, req.Mood,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Insert bundle products
		for _, p := range req.Products {
			_, err := db.DB.Exec(
				"INSERT INTO bundle_products (bundle_id, product_id, quantity) VALUES ($1, $2, $3)",
				req.ID, p.ProductID, p.Quantity,
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Bundle created"})
	})

	// GET ALL BUNDLES
	r.GET("/bundles", func(c *gin.Context) {

		query := `
		SELECT b.id, b.name, b.mood,
		       COALESCE(SUM(p.price * bp.quantity), 0) as total_price
		FROM bundles b
		LEFT JOIN bundle_products bp ON b.id = bp.bundle_id
		LEFT JOIN products p ON p.id = bp.product_id
		GROUP BY b.id;
		`

		rows, err := db.DB.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		type BundleResponse struct {
			ID         string  `json:"id"`
			Name       string  `json:"name"`
			Mood       string  `json:"mood"`
			TotalPrice float64 `json:"total_price"`
		}

		bundles := []BundleResponse{}

		for rows.Next() {
			var b BundleResponse
			if err := rows.Scan(&b.ID, &b.Name, &b.Mood, &b.TotalPrice); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			bundles = append(bundles, b)
		}

		c.JSON(http.StatusOK, bundles)
	})

	r.Run(":8080")
}
