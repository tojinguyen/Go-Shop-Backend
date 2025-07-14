package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/router"
)

//	@title			User Service API
//	@version		1.0
//	@description	User management service for Go-Shop application
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8081
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Description for what is this security definition being used

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Debug configuration (chỉ trong development)
	if cfg.App.IsDevelopment() {
		cfg.Debug()
	}

	// Initialize service container
	serviceContainer, err := container.NewServiceContainer(cfg)
	if err != nil {
		log.Fatal("Failed to initialize services:", err)
	}
	defer serviceContainer.Close()

	// Setup routes with service container
	router := router.SetupRoutes(serviceContainer)

	// Configure HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", serviceContainer.GetConfig().Server.Port),
		Handler:      router,
		ReadTimeout:  serviceContainer.GetConfig().Server.ReadTimeout,
		WriteTimeout: serviceContainer.GetConfig().Server.WriteTimeout,
		IdleTimeout:  serviceContainer.GetConfig().Server.IdleTimeout,
	}

	go startMetricsServer()

	// Start server in a goroutine
	go func() {
		log.Printf("%s starting on %s", serviceContainer.GetConfig().App.Name, serviceContainer.GetConfig().Server.GetServerAddress())
		log.Printf("Environment: %s", serviceContainer.GetConfig().App.Environment)
		log.Printf("Debug Mode: %v", serviceContainer.GetConfig().App.Debug)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a deadline of 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func startMetricsServer() {
	metricsRouter := gin.New()
	metricsRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))
	log.Println("Starting metrics server on :9100")
	if err := metricsRouter.Run(":9100"); err != nil {
		log.Fatalf("Failed to start metrics server: %v", err)
	}
}
