package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/your-username/go-shop/internal/services/user-service/internal/config"
	postgresql_infra "github.com/your-username/go-shop/internal/services/user-service/internal/infra/postgreql-infra"
	redis_infra "github.com/your-username/go-shop/internal/services/user-service/internal/infra/redis-infra"
	"github.com/your-username/go-shop/internal/services/user-service/internal/router"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize PostgreSQL service
	pgService, err := postgresql_infra.NewPostgreSQLService(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to initialize PostgreSQL service:", err)
	}
	defer pgService.Close()

	// Initialize Redis service
	redisService := redis_infra.NewRedisService(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)

	// Test Redis connection
	if err := redisService.Ping(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisService.Close()

	// Setup routes with services
	router := router.SetupRoutes(cfg, pgService, redisService)

	// Configure HTTP server
	server := &http.Server{
		Addr:         cfg.Server.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("%s starting on %s", cfg.App.Name, cfg.Server.GetServerAddress())
		log.Printf("Environment: %s", cfg.App.Environment)
		log.Printf("Debug Mode: %v", cfg.App.Debug)

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
