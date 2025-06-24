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
	port := os.Getenv("USER_SERVICE_PORT")
	if port == "" {
		port = "8000"
	}

	// Create Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "user-service",
			"status":  "healthy",
		})
	})

	// User endpoints
	r.POST("/users/register", registerUser)
	r.POST("/users/login", loginUser)
	r.GET("/users/:id", getUser)
	r.PUT("/users/:id", updateUser)
	r.POST("/users/:id/addresses", addUserAddress)
	r.GET("/users/:id/addresses", getUserAddresses)

	// Start server
	log.Printf("User Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func registerUser(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}

func loginUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   "jwt_token_here",
	})
}

func getUser(c *gin.Context) {
	userID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"email":   "user@example.com",
	})
}

func updateUser(c *gin.Context) {
	userID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"message": "User updated successfully",
	})
}

func addUserAddress(c *gin.Context) {
	userID := c.Param("id")
	c.JSON(http.StatusCreated, gin.H{
		"user_id": userID,
		"message": "Address added successfully",
	})
}

func getUserAddresses(c *gin.Context) {
	userID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"user_id":   userID,
		"addresses": []string{},
	})
}
