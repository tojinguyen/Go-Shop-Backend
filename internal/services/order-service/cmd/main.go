package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/config"
	dependency_container "github.com/toji-dev/go-shop/internal/services/order-service/internal/dependency-container"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/router"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/worker"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return
	}
	fmt.Println("Starting order-service...")

	// Initialize dependency container
	dependencyContainer := dependency_container.NewDependencyContainer(cfg)

	scheduler := worker.NewScheduler(dependencyContainer)
	scheduler.RegisterJobs()

	// Chạy scheduler trong một goroutine riêng
	go scheduler.Start()

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	router.Init(r, dependencyContainer)

	// Start server
	if err := r.Run(cfg.Server.GetServerAddress()); err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}
}
