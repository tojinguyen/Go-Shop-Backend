package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create Gin router
	r := gin.Default()

	// Start server
	log.Printf("Starting shop service on %s", cfg.GetServerAddress())
	if err := r.Run(cfg.GetServerAddress()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
