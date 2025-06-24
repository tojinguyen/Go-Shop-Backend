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
	port := os.Getenv("NOTIFICATION_SERVICE_PORT")
	if port == "" {
		port = "8003"
	}

	// Create Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "notification-service",
			"status":  "healthy",
		})
	})

	// Notification endpoints
	r.POST("/notifications", createNotification)
	r.GET("/notifications/:user_id", getUserNotifications)
	r.PUT("/notifications/:id/read", markAsRead)

	// Start server
	log.Printf("Notification Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func createNotification(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "Notification created successfully",
	})
}

func getUserNotifications(c *gin.Context) {
	userID := c.Param("user_id")
	c.JSON(http.StatusOK, gin.H{
		"user_id":       userID,
		"notifications": []string{},
	})
}

func markAsRead(c *gin.Context) {
	notificationID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"notification_id": notificationID,
		"status":          "read",
	})
}
