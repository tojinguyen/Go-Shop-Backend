package container

import (
	"fmt"

	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/handlers"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/service"
)

// Container holds all the dependencies
type Container struct {
	Config      *config.Config
	DB          *postgresql_infra.PostgreSQLService
	ShopRepo    repository.ShopRepository
	ShopService service.ShopService
	ShopHandler *handlers.ShopHandler
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

	// Initialize services
	shopService := service.NewShopService(shopRepo)

	// Initialize handlers
	shopHandler := handlers.NewShopHandler(shopService)

	return &Container{
		Config:      cfg,
		DB:          db,
		ShopRepo:    shopRepo,
		ShopService: shopService,
		ShopHandler: shopHandler,
	}, nil
}

// Close closes all resources
func (c *Container) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
}
