package main

import (
	"fmt"
	"log"

	"net/http"
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/config"
	dependency_container "github.com/toji-dev/go-shop/internal/services/payment-service/internal/dependency-container"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/router"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/worker"
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
	go scheduler.Start()
	log.Println("Scheduler started.")

	// Initialize Gin router
	r := gin.Default()

	go startMetricsServer()

	go func() {
		log.Println("Starting pprof server on :6065")
		if err := http.ListenAndServe("0.0.0.0:6065", nil); err != nil {
			log.Printf("Pprof server failed to start: %v", err)
		}
	}()

	// Setup routes
	router.Init(r, dependencyContainer)

	// Start server
	if err := r.Run(cfg.Server.GetServerAddress()); err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}
}

func startMetricsServer() {
	metricsRouter := gin.New()
	metricsRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))
	log.Println("Starting metrics server on 0.0.0.0:9100")
	if err := metricsRouter.Run("0.0.0.0:9100"); err != nil {
		log.Fatalf("Failed to start metrics server: %v", err)
	}
}
