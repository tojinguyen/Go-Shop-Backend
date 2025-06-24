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
	port := os.Getenv("ORDER_SERVICE_PORT")
	if port == "" {
		port = "8001"
	}

	// Create Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "order-service",
			"status":  "healthy",
		})
	})

	// Order endpoints
	r.POST("/orders", createOrder)
	r.GET("/orders/:id", getOrder)
	r.GET("/orders/user/:user_id", getUserOrders)
	r.PUT("/orders/:id/status", updateOrderStatus)

	// Start server
	log.Printf("Order Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func createOrder(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
	})
}

func getOrder(c *gin.Context) {
	orderID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"order_id": orderID,
		"status":   "pending",
	})
}

func getUserOrders(c *gin.Context) {
	userID := c.Param("user_id")
	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"orders":  []string{},
	})
}

func updateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"order_id": orderID,
		"message":  "Order status updated successfully",
	})
}
