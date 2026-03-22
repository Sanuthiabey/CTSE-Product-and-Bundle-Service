// @title Product & Bundle Service API
// @version 1.0
// @description API for managing products and bundles
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/docs"

	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/db"
	grpcServer "github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/grpc"
	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/handlers"
	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/middleware"
)

func main() {

	// ==============================
	// LOAD ENV
	// ==============================
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}

	// ==============================
	// DATABASE
	// ==============================
	db.Connect()

	// ==============================
	// START gRPC
	// ==============================
	go grpcServer.StartGRPCServer()

	// ==============================
	// GIN SETUP
	// ==============================
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// ==============================
	// SWAGGER
	// ==============================
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ==============================
	// CORS
	// ==============================
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
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

	// =====================================================
	// 🔓 PUBLIC ROUTES (NO AUTH)
	// =====================================================
	public := r.Group("/api")
	{
		public.GET("/products", handlers.GetProducts)
		public.GET("/products/:id", handlers.GetProduct)
		public.GET("/bundles", handlers.GetBundles)
	}

	// =====================================================
	// 🔐 AUTH REQUIRED ROUTES
	// =====================================================
	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired())
	{
		protected.POST("/stock/validate", handlers.ValidateStock)
	}

	// =====================================================
	// 👑 ADMIN ROUTES
	// =====================================================
	admin := r.Group("/admin")
	admin.Use(middleware.AuthRequired())
	admin.Use(middleware.AdminOnly())
	{
		// PRODUCTS
		admin.POST("/products", handlers.CreateProduct)
		admin.PUT("/products/:id", handlers.UpdateProduct)
		admin.DELETE("/products/:id", handlers.DeleteProduct)

		// BUNDLES
		admin.POST("/bundles", handlers.CreateBundle)

		// STOCK
		admin.POST("/stock/reduce", handlers.ReduceStock)
		admin.PUT("/stock/update", handlers.SetStock)
	}
	// START SERVER
	r.Run(":8080")
}
