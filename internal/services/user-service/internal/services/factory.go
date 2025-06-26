package services

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/your-username/go-shop/internal/services/user-service/internal/handlers"
	"github.com/your-username/go-shop/internal/services/user-service/internal/middleware"
)

// HandlerFactory creates handlers with all necessary dependencies
type HandlerFactory struct {
	container *ServiceContainer
}

// NewHandlerFactory creates a new handler factory
func NewHandlerFactory(container *ServiceContainer) *HandlerFactory {
	return &HandlerFactory{
		container: container,
	}
}

// CreateAuthHandler creates an authentication handler
func (hf *HandlerFactory) CreateAuthHandler() *handlers.AuthHandler {
	return handlers.NewAuthHandler(
		hf.container.GetJWT(),
		hf.container.GetConfig(),
		hf.container.GetPostgreSQL(),
		hf.container.GetRedis(),
	)
}

// CreateProfileHandler creates a profile handler
func (hf *HandlerFactory) CreateProfileHandler() *handlers.ProfileHandler {
	return handlers.NewProfileHandler(
		hf.container.GetJWT(),
		hf.container.GetConfig(),
		hf.container.GetPostgreSQL(),
		hf.container.GetRedis(),
	)
}

// CreateAuthMiddleware creates authentication middleware
func (hf *HandlerFactory) CreateAuthMiddleware() func(c *gin.Context) {
	return middleware.AuthMiddleware(hf.container.GetJWT())
}

// HealthChecker provides health check functionality
type HealthChecker struct {
	container *ServiceContainer
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(container *ServiceContainer) *HealthChecker {
	return &HealthChecker{
		container: container,
	}
}

// CheckHealth performs health checks on all services
func (hc *HealthChecker) CheckHealth(ctx context.Context) map[string]interface{} {
	health := make(map[string]interface{})

	// App info
	health["service"] = "user-service"
	health["version"] = hc.container.GetConfig().App.Version
	health["environment"] = hc.container.GetConfig().App.Environment
	health["timestamp"] = time.Now().UTC()

	// Check PostgreSQL
	pgHealth := "healthy"
	if hc.container.PostgreSQL != nil {
		if conn, err := hc.container.PostgreSQL.GetConnection(ctx); err != nil {
			pgHealth = "unhealthy: " + err.Error()
		} else {
			conn.Release()
		}
	} else {
		pgHealth = "not initialized"
	}
	health["postgresql"] = pgHealth

	// Check Redis
	redisHealth := "healthy"
	if hc.container.Redis != nil {
		if err := hc.container.Redis.Ping(); err != nil {
			redisHealth = "unhealthy: " + err.Error()
		}
	} else {
		redisHealth = "not initialized"
	}
	health["redis"] = redisHealth

	// Overall status
	if pgHealth == "healthy" && redisHealth == "healthy" {
		health["status"] = "healthy"
	} else {
		health["status"] = "degraded"
	}

	return health
}

// ServiceManager manages service lifecycle
type ServiceManager struct {
	container *ServiceContainer
}

// NewServiceManager creates a new service manager
func NewServiceManager(container *ServiceContainer) *ServiceManager {
	return &ServiceManager{
		container: container,
	}
}

// Restart restarts all services (useful for development)
func (sm *ServiceManager) Restart() error {
	// Close existing services
	sm.container.Close()

	// Reinitialize services
	if err := sm.container.initPostgreSQL(); err != nil {
		return err
	}

	if err := sm.container.initRedis(); err != nil {
		return err
	}

	sm.container.initJWT()

	return nil
}

// GetMetrics returns service metrics
func (sm *ServiceManager) GetMetrics() map[string]interface{} {
	metrics := make(map[string]interface{})

	// PostgreSQL metrics
	if sm.container.PostgreSQL != nil {
		pool := sm.container.PostgreSQL.GetPool()
		if pool != nil {
			stat := pool.Stat()
			metrics["postgresql"] = map[string]interface{}{
				"total_connections":    stat.TotalConns(),
				"acquired_connections": stat.AcquiredConns(),
				"idle_connections":     stat.IdleConns(),
				"max_connections":      stat.MaxConns(),
			}
		}
	}

	// Add more metrics as needed
	metrics["uptime"] = time.Since(time.Now()).String() // This would be tracked from startup

	return metrics
}
