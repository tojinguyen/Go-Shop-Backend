package services

import (
	"fmt"
	"log"

	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	redis_infra "github.com/toji-dev/go-shop/internal/pkg/infra/redis-infra"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/config"
	jwtService "github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/jwt"
)

// ServiceContainer holds all application services
type ServiceContainer struct {
	Config     *config.Config
	PostgreSQL *postgresql_infra.PostgreSQLService
	Redis      *redis_infra.RedisService
	JWT        jwtService.JwtService
}

// NewServiceContainer creates and initializes all services
func NewServiceContainer(cfg *config.Config) (*ServiceContainer, error) {
	container := &ServiceContainer{
		Config: cfg,
	}

	// Initialize PostgreSQL service
	if err := container.initPostgreSQL(); err != nil {
		return nil, fmt.Errorf("failed to initialize PostgreSQL: %w", err)
	}

	// Initialize Redis service
	if err := container.initRedis(); err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	// Initialize JWT service
	container.initJWT()

	log.Println("All services initialized successfully")
	return container, nil
}

// initPostgreSQL initializes PostgreSQL service
func (sc *ServiceContainer) initPostgreSQL() error {
	// Convert config types
	pgConfig := &postgresql_infra.DatabaseConfig{
		Host:         sc.Config.Database.Host,
		Port:         sc.Config.Database.Port,
		User:         sc.Config.Database.User,
		Password:     sc.Config.Database.Password,
		Name:         sc.Config.Database.Name,
		SSLMode:      sc.Config.Database.SSLMode,
		MaxOpenConns: sc.Config.Database.MaxOpenConns,
		MaxIdleConns: sc.Config.Database.MaxIdleConns,
		MaxLifetime:  sc.Config.Database.MaxLifetime,
	}

	pgService, err := postgresql_infra.NewPostgreSQLService(pgConfig)
	if err != nil {
		return fmt.Errorf("failed to create PostgreSQL service: %w", err)
	}

	sc.PostgreSQL = pgService
	log.Println("PostgreSQL service initialized")
	return nil
}

// initRedis initializes Redis service
func (sc *ServiceContainer) initRedis() error {
	redisService := redis_infra.NewRedisService(sc.Config.Redis.Host, sc.Config.Redis.Port, sc.Config.Redis.Password, sc.Config.Redis.DB)

	// Test Redis connection
	if err := redisService.Ping(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	sc.Redis = redisService
	log.Println("Redis service initialized")
	return nil
}

// initJWT initializes JWT service
func (sc *ServiceContainer) initJWT() {
	sc.JWT = jwtService.NewJwtService(sc.Config)
	log.Println("JWT service initialized")
}

// Close gracefully closes all services
func (sc *ServiceContainer) Close() {
	log.Println("Shutting down services...")

	if sc.PostgreSQL != nil {
		sc.PostgreSQL.Close()
		log.Println("PostgreSQL service closed")
	}

	if sc.Redis != nil {
		sc.Redis.Close()
		log.Println("Redis service closed")
	}

	log.Println("All services closed")
}

// GetPostgreSQL returns PostgreSQL service
func (sc *ServiceContainer) GetPostgreSQL() *postgresql_infra.PostgreSQLService {
	return sc.PostgreSQL
}

// GetRedis returns Redis service
func (sc *ServiceContainer) GetRedis() *redis_infra.RedisService {
	return sc.Redis
}

// GetJWT returns JWT service
func (sc *ServiceContainer) GetJWT() jwtService.JwtService {
	return sc.JWT
}

// GetConfig returns configuration
func (sc *ServiceContainer) GetConfig() *config.Config {
	return sc.Config
}
