package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/router"
)

func main() {
	// Load configuration
	_, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return
	}
	fmt.Println("Starting order-service...")

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	api := r.Group("/api")
	router.Init(api)

	// Start server
	if err := r.Run(":8082"); err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}
}
