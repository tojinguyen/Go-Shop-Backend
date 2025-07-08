package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/config"
	dependency_container "github.com/toji-dev/go-shop/internal/services/product-service/internal/dependency-container"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/router"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	dependencyContainer, err := dependency_container.NewDependencyContainer(cfg)

	if err != nil {
		log.Fatalf("Failed to initialize dependency container: %v", err)
	}
	dependencyContainer.Close()

	router := router.SetupRoutes(r, dependencyContainer)

	server := &http.Server{
		Addr:         dependencyContainer.GetConfig().Server.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  dependencyContainer.GetConfig().Server.ReadTimeout,
		WriteTimeout: dependencyContainer.GetConfig().Server.WriteTimeout,
		IdleTimeout:  dependencyContainer.GetConfig().Server.IdleTimeout,
	}

	go func() {
		log.Printf("%s starting on %s", dependencyContainer.GetConfig().App.Name, dependencyContainer.GetConfig().Server.GetServerAddress())
		log.Printf("Environment: %s", dependencyContainer.GetConfig().App.Environment)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
}
