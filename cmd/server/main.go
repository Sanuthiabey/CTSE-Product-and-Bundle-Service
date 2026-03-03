package main

import (
	"net/http"

	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/models"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	// In-memory storage
	products := make(map[string]models.Product)

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "product-and-bundle-service",
			"status":  "running",
		})
	})

	// CREATE product
	r.POST("/products", func(c *gin.Context) {
		var product models.Product

		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		products[product.ID] = product
		c.JSON(http.StatusCreated, product)
	})

	// GET all products
	r.GET("/products", func(c *gin.Context) {
		list := []models.Product{}

		for _, p := range products {
			list = append(list, p)
		}

		c.JSON(http.StatusOK, list)
	})

	// GET product by ID
	r.GET("/products/:id", func(c *gin.Context) {
		id := c.Param("id")

		product, exists := products[id]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusOK, product)
	})

	// UPDATE product
	r.PUT("/products/:id", func(c *gin.Context) {
		id := c.Param("id")

		_, exists := products[id]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		var updated models.Product
		if err := c.ShouldBindJSON(&updated); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updated.ID = id
		products[id] = updated

		c.JSON(http.StatusOK, updated)
	})

	// DELETE product
	r.DELETE("/products/:id", func(c *gin.Context) {
		id := c.Param("id")

		_, exists := products[id]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		delete(products, id)
		c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
	})

	r.Run(":8080")
}
