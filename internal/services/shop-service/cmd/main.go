package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get port from environment or use default
	port := os.Getenv("SHOP_SERVICE_PORT")
	if port == "" {
		port = "8006"
	}

	// Create Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "shop-service",
			"status":  "healthy",
		})
	})

	// Shop endpoints
	r.POST("/shops", createShop)
	r.GET("/shops/:id", getShop)
	r.PUT("/shops/:id", updateShop)
	r.GET("/shops", getShops)
	r.POST("/shops/:id/promotions", createPromotion)
	r.GET("/shops/:id/promotions", getShopPromotions)

	// Start server
	log.Printf("Shop Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func createShop(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "Shop created successfully",
	})
}

func getShop(c *gin.Context) {
	shopID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"shop_id": shopID,
		"name":    "Sample Shop",
	})
}

func updateShop(c *gin.Context) {
	shopID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"shop_id": shopID,
		"message": "Shop updated successfully",
	})
}

func getShops(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"shops": []string{},
	})
}

func createPromotion(c *gin.Context) {
	shopID := c.Param("id")
	c.JSON(http.StatusCreated, gin.H{
		"shop_id": shopID,
		"message": "Promotion created successfully",
	})
}

func getShopPromotions(c *gin.Context) {
	shopID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"shop_id":    shopID,
		"promotions": []string{},
	})
}
