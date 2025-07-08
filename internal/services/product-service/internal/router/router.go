package router

import (
	"github.com/gin-gonic/gin"
	dependency_container "github.com/toji-dev/go-shop/internal/services/product-service/internal/dependency-container"
)

func SetupRoutes(r *gin.Engine, serviceContainer *dependency_container.DependencyContainer) *gin.Engine {
	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	return r
}
