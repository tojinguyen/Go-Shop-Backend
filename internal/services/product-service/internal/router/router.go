package router

import (
	"github.com/gin-gonic/gin"
	dependency_container "github.com/toji-dev/go-shop/internal/services/product-service/internal/dependency-container"
)

func SetupRoutes(r *gin.Engine, serviceContainer *dependency_container.DependencyContainer) *gin.Engine {
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})

	return r
}
