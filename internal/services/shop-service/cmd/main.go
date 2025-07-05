package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/api"
	createshop "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/commands/create_shop"
	deleteshop "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/commands/delete_shop"
	updateshop "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/commands/update_shop"
	getshop "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/queries/get_shop"
	getshops "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/queries/get_shops"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	dbConfig := &postgresql_infra.DatabaseConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	db, err := postgresql_infra.NewPostgreSQLService(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	shopRepo := repository.NewPostgresShopRepository(db)

	// Set Gin mode based on environment
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.Default()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "shop-service",
			"version": cfg.App.Version,
		})
	})

	// Initialize feature handlers
	// Create shop
	createShopHandler := createshop.NewHandler(shopRepo)
	createShopAPIHandler := createshop.NewAPIHandler(createShopHandler)

	// Get shop by ID
	getShopHandler := getshop.NewHandler(shopRepo)
	getShopAPIHandler := getshop.NewAPIHandler(getShopHandler)

	// Get shops by owner
	getShopsHandler := getshops.NewQueryHandler(shopRepo)
	getShopsAPIHandler := getshops.NewAPIHandler(getShopsHandler)

	// Update shop
	updateShopHandler := updateshop.NewCommandHandler(shopRepo)
	updateShopAPIHandler := updateshop.NewAPIHandler(updateShopHandler)

	// Delete shop
	deleteShopHandler := deleteshop.NewCommandHandler(shopRepo)
	deleteShopAPIHandler := deleteshop.NewAPIHandler(deleteShopHandler)

	// Register shop routes
	api.RegisterShopRoutes(
		r,
		createShopAPIHandler,
		getShopAPIHandler,
		getShopsAPIHandler,
		updateShopAPIHandler,
		deleteShopAPIHandler,
	)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting shop service on %s", cfg.GetServerAddress())
		if err := r.Run(cfg.GetServerAddress()); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down shop service...")
}
