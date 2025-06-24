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
	port := os.Getenv("PAYMENT_SERVICE_PORT")
	if port == "" {
		port = "8002"
	}

	// Create Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "payment-service",
			"status":  "healthy",
		})
	})

	// Payment endpoints
	r.POST("/payments", processPayment)
	r.GET("/payments/:id", getPayment)
	r.POST("/payments/:id/refund", refundPayment)

	// Start server
	log.Printf("Payment Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func processPayment(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "Payment processed successfully",
	})
}

func getPayment(c *gin.Context) {
	paymentID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"payment_id": paymentID,
		"status":     "completed",
	})
}

func refundPayment(c *gin.Context) {
	paymentID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"payment_id": paymentID,
		"message":    "Refund processed successfully",
	})
}
