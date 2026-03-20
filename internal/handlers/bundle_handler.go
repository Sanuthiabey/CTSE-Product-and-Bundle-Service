package handlers

import (
	"net/http"

	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/models"
	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/services"
	"github.com/gin-gonic/gin"
)

// CreateBundle godoc
// @Summary Create bundle
// @Tags Bundles
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param Role header string true "Admin role"
// @Security BearerAuth
// @Param bundle body models.CreateBundleRequest true "Bundle"
// @Success 201 {object} models.Bundle
// @Router /admin/bundles [post]
func CreateBundle(c *gin.Context) {
	var req models.CreateBundleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bundle, err := services.CreateBundle(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bundle)
}

// GetBundles godoc
// @Summary Get all bundles
// @Tags Bundles
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Security BearerAuth
// @Success 200 {array} models.Bundle
// @Router /api/bundles [get]
func GetBundles(c *gin.Context) {

	bundles, err := services.GetBundles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bundles)
}
