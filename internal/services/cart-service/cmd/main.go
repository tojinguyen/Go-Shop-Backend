package main

import (
	"log"

	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/config"
	dependency_container "github.com/toji-dev/go-shop/internal/services/cart-service/internal/dependency_container"
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
}
