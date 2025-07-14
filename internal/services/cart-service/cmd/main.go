package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/config"
	dependency_container "github.com/toji-dev/go-shop/internal/services/cart-service/internal/dependency-container"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	log.Printf("Successfully loaded configuration for %s", cfg.App.Name)
	log.Printf("Server will run on port: %s", cfg.Server.Port)
	log.Printf("Connecting to Redis on: %s:%s (DB: %d)", cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.DB)

	dependencyContainer, err := dependency_container.NewDependencyContainer(cfg)

	if err != nil {
		log.Fatal("Error initializing dependency container:", err)
	}

	defer dependencyContainer.Close()

	g := gin.New()
	g.Use(gin.Recovery())
	g.Use(gin.Logger())

	router.SetupRoutes(g, dependencyContainer)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("%s starting on %s", dependencyContainer.GetConfig().App.Name, dependencyContainer.GetConfig().Server.GetServerAddress())
		log.Printf("Environment: %s", dependencyContainer.GetConfig().App.Environment)

		if err := g.Run(dependencyContainer.GetConfig().Server.GetServerAddress()); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down shop service...")
}
