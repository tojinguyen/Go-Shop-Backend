package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/toji-dev/go-shop/internal/services/user-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/router"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize service container
	serviceContainer, err := services.NewServiceContainer(cfg)
	if err != nil {
		log.Fatal("Failed to initialize services:", err)
	}
	defer serviceContainer.Close()

	// Setup routes with service container
	router := router.SetupRoutes(serviceContainer)

	// Configure HTTP server
	server := &http.Server{
		Addr:         serviceContainer.GetConfig().Server.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  serviceContainer.GetConfig().Server.ReadTimeout,
		WriteTimeout: serviceContainer.GetConfig().Server.WriteTimeout,
		IdleTimeout:  serviceContainer.GetConfig().Server.IdleTimeout,
	}

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
