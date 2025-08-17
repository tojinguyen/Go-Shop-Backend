package dependency_container

import (
	"fmt"
	"log"

	redis_infra "github.com/toji-dev/go-shop/internal/pkg/infra/redis-infra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/pkg/jwt"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/config"
	grpc "github.com/toji-dev/go-shop/internal/services/cart-service/internal/grpc/adapter"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/repository"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/usecase"
)

type DependencyContainer struct {
	config          *config.Config
	postgreSQL      *postgresql_infra.PostgreSQLService
	redis           *redis_infra.RedisService
	gormDB          *gorm.DB
	cartRepository  repository.CartRepository
	cartUseCase     usecase.CartUseCase
	cartItemUseCase usecase.CartItemUseCase
	product_adapter grpc.ProductServiceAdapter
	jwt             jwt.JwtService
}

func (sc *DependencyContainer) GetConfig() *config.Config {
	return sc.config
}

func (sc *DependencyContainer) GetRedisService() *redis_infra.RedisService {
	return sc.redis
}

func (sc *DependencyContainer) GetCartUseCase() usecase.CartUseCase {
	return sc.cartUseCase
}

func (sc *DependencyContainer) GetCartItemUseCase() usecase.CartItemUseCase {
	return sc.cartItemUseCase
}

func NewDependencyContainer(cfg *config.Config) (*DependencyContainer, error) {
	container := &DependencyContainer{
		config: cfg,
	}

	// Initialize PostgreSQL service
	if err := container.initPostgreSQL(); err != nil {
		return nil, err
	}

	// Initialize GORM
	if err := container.initGORM(); err != nil {
		return nil, err
	}

	// Initialize Redis service
	if err := container.initRedis(); err != nil {
		return nil, err
	}

	// Initialize repositories
	container.initRepositories()

	// Initialize gRPC service adapter
	container.initGrpcServiceAdapter()

	// Initialize use cases
	container.initUseCases()

	// Initialize JWT service
	container.initJwtService()

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

func (sc *DependencyContainer) initGORM() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		sc.config.Database.Host, sc.config.Database.User, sc.config.Database.Password, sc.config.Database.DBName, sc.config.Database.Port, sc.config.Database.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	sc.gormDB = db
	log.Println("GORM initialized")
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

func (sc *DependencyContainer) initRepositories() {
	sc.cartRepository = repository.NewCartRepository(sc.postgreSQL)
	log.Println("Repositories initialized")
}

func (sc *DependencyContainer) initUseCases() {
	sc.cartUseCase = usecase.NewCartUseCase(sc.cartRepository)
	sc.cartItemUseCase = usecase.NewCartItemUseCase(sc.cartRepository, sc.product_adapter)
	log.Println("Use cases initialized")
}

func (sc *DependencyContainer) initGrpcServiceAdapter() error {
	if sc.product_adapter != nil {
		return nil
	}

	productServiceAddr := fmt.Sprintf("%s:%d", sc.config.Grpc.ProductServiceHost, sc.config.Grpc.ProductServicePort)
	adapter, err := grpc.NewGrpcProductAdapter(productServiceAddr)
	if err != nil {
		return fmt.Errorf("failed to create product service adapter: %w", err)
	}
	sc.product_adapter = adapter
	log.Println("Product service adapter initialized")
	return nil
}

func (sc *DependencyContainer) initJwtService() error {
	jwtCfg := jwt.JWTConfig{
		SecretKey:       sc.config.Jwt.SecretKey,
		AccessTokenTTL:  sc.config.Jwt.AccessTokenTTL,
		RefreshTokenTTL: sc.config.Jwt.RefreshTokenTTL,
		Issuer:          sc.config.Jwt.Issuer,
	}
	sc.jwt = jwt.NewJwtService(jwtCfg)
	log.Println("JWT service initialized")
	return nil
}

func (sc *DependencyContainer) GetJwtService() jwt.JwtService {
	if sc.jwt == nil {
		if err := sc.initJwtService(); err != nil {
			return nil
		}
	}
	return sc.jwt
}

func (sc *DependencyContainer) GetProductServiceAdapter() (grpc.ProductServiceAdapter, error) {
	if sc.product_adapter == nil {
		if err := sc.initGrpcServiceAdapter(); err != nil {
			return nil, err
		}
	}
	return sc.product_adapter, nil
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
}
