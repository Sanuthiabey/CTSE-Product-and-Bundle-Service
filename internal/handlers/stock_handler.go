package handlers

import (
	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/models"
	"net/http"

	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/services"
	"github.com/gin-gonic/gin"
)

// ValidateStock godoc
// @Summary Validate stock
// @Tags Stock
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Security BearerAuth
// @Param request body map[string]interface{} true "Stock request"
// @Success 200 {object} map[string]interface{}
// @Router /api/stock/validate [post]
func ValidateStock(c *gin.Context) {

	var req struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.ValidateStock(req.ProductID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ReduceStock godoc
// @Summary Reduce stock
// @Tags Stock
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param Role header string true "Admin role"
// @Security BearerAuth
// @Param request body map[string]interface{} true "Stock reduction"
// @Success 200 {object} map[string]interface{}
// @Router /admin/stock/reduce [post]
func ReduceStock(c *gin.Context) {

	var req struct {
		Items []models.StockItem `json:"items"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.ReduceStock(req.Items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stock reduced",
	})
}
