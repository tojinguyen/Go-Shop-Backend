package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/go-shop/internal/services/user-service/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Set Gin mode based on environment
	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": cfg.App.Name,
			"version": cfg.App.Version,
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
	serverAddr := cfg.Server.GetServerAddress()
	log.Printf("%s starting on %s", cfg.App.Name, serverAddr)
	if err := r.Run(serverAddr); err != nil {
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
