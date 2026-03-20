// @title Product & Bundle Service API
// @version 1.0
// @description API for managing products and bundles
// @host localhost:8080
// @BasePath /

package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"

	_ "github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/docs"

	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/db"
	grpcServer "github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/grpc"
	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/handlers"
	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/middleware"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("No .env file found")
	}
	// DATABASE CONNECTION
	db.Connect()

	// START gRPC SERVER

	go grpcServer.StartGRPCServer()

	// GIN SERVER
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// CORS CONFIG
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// HEALTH CHECK
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "product-and-bundle-service",
			"status":  "running",
		})
	})

	// PUBLIC ROUTES
	api := r.Group("/api")
	api.Use(middleware.AuthRequired())

	{
		// All logged users
		api.GET("/products", handlers.GetProducts)
		api.GET("/products/:id", handlers.GetProduct)
	}

	// ADMIN ROUTES
	admin := r.Group("/admin")
	admin.Use(middleware.AuthRequired())
	admin.Use(middleware.AdminOnly())

	{
		admin.POST("/products", handlers.CreateProduct)
		admin.PUT("/products/:id", handlers.UpdateProduct)
		admin.DELETE("/products/:id", handlers.DeleteProduct)
	}

	// START HTTP SERVER
	r.Run(":8080")
}
