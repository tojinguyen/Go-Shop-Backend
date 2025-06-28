package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/middleware"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/services"
)

// HandlerFactory creates handlers with all necessary dependencies
type HandlerFactory struct {
	container *services.ServiceContainer
}

// NewHandlerFactory creates a new handler factory
func NewHandlerFactory(container *services.ServiceContainer) *HandlerFactory {
	return &HandlerFactory{
		container: container,
	}
}

// CreateAuthHandler creates an authentication handler
func (hf *HandlerFactory) CreateAuthHandler() *AuthHandler {
	return NewAuthHandler(
		hf.container.GetJWT(),
		hf.container.GetConfig(),
		hf.container.GetPostgreSQL(),
		hf.container.GetRedis(),
	)
}

// CreateProfileHandler creates a profile handler
func (hf *HandlerFactory) CreateProfileHandler() *ProfileHandler {
	return NewProfileHandler(
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
