package container

import (
	"fmt"
	"log"

	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	redis_infra "github.com/toji-dev/go-shop/internal/pkg/infra/redis-infra"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/config"
	jwtService "github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/jwt"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/repository"
)

// ServiceContainer holds all application services
type ServiceContainer struct {
	config          *config.Config
	postgreSQL      *postgresql_infra.PostgreSQLService
	redis           *redis_infra.RedisService
	jwt             jwtService.JwtService
	userAccountRepo repository.UserAccountRepository
}

// NewServiceContainer creates and initializes all services
func NewServiceContainer(cfg *config.Config) (ServiceContainer, error) {
	container := ServiceContainer{
		config: cfg,
	}

	// Initialize PostgreSQL service
	if err := container.initPostgreSQL(); err != nil {
		return ServiceContainer{}, fmt.Errorf("failed to initialize PostgreSQL: %w", err)
	}

	// Initialize Redis service
	if err := container.initRedis(); err != nil {
		return ServiceContainer{}, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	// Initialize JWT service
	container.initJWT()

	// Initialize UserAccountRepository
	container.initUserAccountRepository()

	log.Println("All services initialized successfully")
	return container, nil
}

// initPostgreSQL initializes PostgreSQL service
func (sc *ServiceContainer) initPostgreSQL() error {
	// Convert config types
	pgConfig := &postgresql_infra.DatabaseConfig{
		Host:         sc.config.Database.Host,
		Port:         sc.config.Database.Port,
		User:         sc.config.Database.User,
		Password:     sc.config.Database.Password,
		Name:         sc.config.Database.Name,
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
func (sc *ServiceContainer) initRedis() error {
	redisService := redis_infra.NewRedisService(sc.config.Redis.Host, sc.config.Redis.Port, sc.config.Redis.Password, sc.config.Redis.DB)

	// Test Redis connection
	if err := redisService.Ping(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	sc.redis = redisService
	log.Println("Redis service initialized")
	return nil
}

// initJWT initializes JWT service
func (sc *ServiceContainer) initJWT() {
	sc.jwt = jwtService.NewJwtService(sc.config)
	log.Println("JWT service initialized")
}

func (sc *ServiceContainer) initUserAccountRepository() {
	sc.userAccountRepo = repository.NewUserAccountRepository(sc.postgreSQL)
	log.Println("UserAccountRepository initialized")
}

// Close gracefully closes all services
func (sc *ServiceContainer) Close() {
	log.Println("Shutting down services...")

	if sc.postgreSQL != nil {
		sc.postgreSQL.Close()
		log.Println("PostgreSQL service closed")
	}

	if sc.redis != nil {
		sc.redis.Close()
		log.Println("Redis service closed")
	}

	log.Println("All services closed")
}

// GetPostgreSQL returns PostgreSQL service
func (sc *ServiceContainer) GetPostgreSQL() *postgresql_infra.PostgreSQLService {
	return sc.postgreSQL
}

// GetRedis returns Redis service
func (sc *ServiceContainer) GetRedis() *redis_infra.RedisService {
	return sc.redis
}

// GetJWT returns JWT service
func (sc *ServiceContainer) GetJWT() jwtService.JwtService {
	return sc.jwt
}

// GetConfig returns configuration
func (sc *ServiceContainer) GetConfig() *config.Config {
	return sc.config
}

// GetUserAccountRepo returns UserAccountRepository
func (sc *ServiceContainer) GetUserAccountRepo() repository.UserAccountRepository {
	return sc.userAccountRepo
}
