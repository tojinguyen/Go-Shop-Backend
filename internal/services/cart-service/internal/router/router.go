package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dependency_container "github.com/toji-dev/go-shop/internal/services/cart-service/internal/dependency-container"
)

func SetupRoutes(r *gin.Engine, dependencyContainer *dependency_container.DependencyContainer) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "cart-service",
		})
	})
}
