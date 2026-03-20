package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/models"
	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/services"
)

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product"
// @Success 201 {object} models.Product
// @Failure 400 {object} map[string]string
// @Router /admin/products [post]
func CreateProduct(c *gin.Context) {

	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.CreateProduct(product)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProducts godoc
// @Summary Get all products
// @Description Retrieve list of products
// @Tags Products
// @Produce json
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {array} models.Product
// @Router /api/products [get]
func GetProducts(c *gin.Context) {

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	products, err := services.GetAllProducts(limit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct godoc
// @Summary Get product by ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 404 {object} map[string]string
// @Router /api/products/{id} [get]
func GetProduct(c *gin.Context) {

	id := c.Param("id")

	product, err := services.GetProductByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary Update product
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.Product true "Updated Product"
// @Success 200 {object} models.Product
// @Router /admin/products/{id} [put]
func UpdateProduct(c *gin.Context) {

	id := c.Param("id")

	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.UpdateProduct(id, product)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	product.ID = id

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete product
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]string
// @Router /admin/products/{id} [delete]
func DeleteProduct(c *gin.Context) {

	id := c.Param("id")

	err := services.DeleteProduct(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted",
	})
}
