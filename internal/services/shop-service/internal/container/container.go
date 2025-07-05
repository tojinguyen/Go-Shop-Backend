package container

import (
	"fmt"

	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/config"
	createshop "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/commands/create_shop"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository"
)

// Container holds all the dependencies
type Container struct {
	Config               *config.Config
	DB                   *postgresql_infra.PostgreSQLService
	ShopRepo             repository.ShopRepository
	CreateShopHandler    *createshop.Handler
	CreateShopAPIHandler *createshop.APIHandler
}

// NewContainer creates a new container with all dependencies
func NewContainer(cfg *config.Config) (*Container, error) {
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
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize repositories
	shopRepo := repository.NewPostgresShopRepository(db)

	// Initialize feature handlers
	createShopHandler := createshop.NewHandler(shopRepo)
	createShopAPIHandler := createshop.NewAPIHandler(createShopHandler)

	return &Container{
		Config:               cfg,
		DB:                   db,
		ShopRepo:             shopRepo,
		CreateShopHandler:    createShopHandler,
		CreateShopAPIHandler: createShopAPIHandler,
	}, nil
}

// Close closes all resources
func (c *Container) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
}
