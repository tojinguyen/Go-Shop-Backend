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
	port := os.Getenv("REVIEW_SERVICE_PORT")
	if port == "" {
		port = "8004"
	}

	// Create Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "review-service",
			"status":  "healthy",
		})
	})

	// Review endpoints
	r.POST("/reviews", createReview)
	r.GET("/reviews/product/:product_id", getProductReviews)
	r.GET("/reviews/:id", getReview)
	r.PUT("/reviews/:id", updateReview)
	r.DELETE("/reviews/:id", deleteReview)

	// Start server
	log.Printf("Review Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func createReview(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "Review created successfully",
	})
}

func getProductReviews(c *gin.Context) {
	productID := c.Param("product_id")
	c.JSON(http.StatusOK, gin.H{
		"product_id": productID,
		"reviews":    []string{},
	})
}

func getReview(c *gin.Context) {
	reviewID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"review_id": reviewID,
		"rating":    5,
	})
}

func updateReview(c *gin.Context) {
	reviewID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"review_id": reviewID,
		"message":   "Review updated successfully",
	})
}

func deleteReview(c *gin.Context) {
	reviewID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"review_id": reviewID,
		"message":   "Review deleted successfully",
	})
}
