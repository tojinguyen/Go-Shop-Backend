package dependency_container

import (
	"fmt"
	"log"

	redis_infra "github.com/toji-dev/go-shop/internal/pkg/infra/redis-infra"

	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/repository"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/service"
)

type DependencyContainer struct {
	config      *config.Config
	postgreSQL  *postgresql_infra.PostgreSQLService
	redis       *redis_infra.RedisService
	productRepo repository.ProductRepository
	shopService service.ShopServiceAdapter
}

func (sc *DependencyContainer) GetConfig() *config.Config {
	return sc.config
}

func (sc *DependencyContainer) GetRedisService() *redis_infra.RedisService {
	return sc.redis
}

func (sc *DependencyContainer) GetProductRepository() repository.ProductRepository {
	return sc.productRepo
}

func (sc *DependencyContainer) GetShopServiceAdapter() service.ShopServiceAdapter {
	return sc.shopService
}

func NewDependencyContainer(cfg *config.Config) (*DependencyContainer, error) {
	container := &DependencyContainer{
		config: cfg,
	}

	// Initialize PostgreSQL service
	if err := container.initPostgreSQL(); err != nil {
		return nil, err
	}

	// Initialize Redis service
	if err := container.initRedis(); err != nil {
		return nil, err
	}

	// Initialize repositories
	container.initProductRepository()

	// Initialize shop service adapter
	if err := container.initShopServiceAdapter(); err != nil {
		return nil, fmt.Errorf("failed to initialize shop service adapter: %w", err)
	}

	return container, nil
}

func (sc *DependencyContainer) initPostgreSQL() error {
	// Convert config types
	pgConfig := &postgresql_infra.DatabaseConfig{
		Host:         sc.config.Database.Host,
		Port:         sc.config.Database.Port,
		User:         sc.config.Database.User,
		Password:     sc.config.Database.Password,
		Name:         sc.config.Database.DBName,
		SSLMode:      sc.config.Database.SSLMode,
		MaxOpenConns: sc.config.Database.MaxOpenConns,
		MaxIdleConns: sc.config.Database.MaxIdleConns,
		MaxLifetime:  sc.config.Database.MaxLifetime,
	}

	pgService, err := postgresql_infra.NewPostgreSQLService(pgConfig)
	if err != nil {
		return fmt.Errorf("failed to create PostgreSQL service: %w", err)
	}

	sc.postgreSQL = pgService
	log.Println("PostgreSQL service initialized")
	return nil
}

// initRedis initializes Redis service
func (sc *DependencyContainer) initRedis() error {
	redisService := redis_infra.NewRedisService(sc.config.Redis.Host, sc.config.Redis.Port, sc.config.Redis.Password, sc.config.Redis.DB)

	// Test Redis connection
	if err := redisService.Ping(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	sc.redis = redisService
	log.Println("Redis service initialized")
	return nil
}

func (sc *DependencyContainer) initProductRepository() {
	sc.productRepo = repository.NewProductRepository(sc.postgreSQL)
	log.Println("Product repository initialized")
}

func (sc *DependencyContainer) initShopServiceAdapter() error {
	shopService, err := service.NewGrpcShopAdapter(sc.config.ShopService.Address)
	if err != nil {
		return fmt.Errorf("failed to create shop service adapter: %w", err)
	}

	sc.shopService = shopService
	log.Println("Shop service adapter initialized")
	return nil
}

func (sc *DependencyContainer) Close() {
	if sc.postgreSQL != nil {
		sc.postgreSQL.Close()
		log.Println("PostgreSQL service closed")
	}

	if sc.redis != nil {
		sc.redis.Close()
		log.Println("Redis service closed")
	}

	if sc.shopService != nil {
		sc.shopService.Close()
		log.Println("Shop service adapter closed")
	}
}
