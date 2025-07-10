package router

import (
	"github.com/gin-gonic/gin"
	dependency_container "github.com/toji-dev/go-shop/internal/services/product-service/internal/dependency-container"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/handler"
)

func SetupRoutes(r *gin.Engine, serviceContainer *dependency_container.DependencyContainer) *gin.Engine {
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	productHandler := handler.NewProductHandler(serviceContainer.GetProductRepository(), serviceContainer.GetRedisService())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})

	v1 := r.Group("/api/v1")
	{
		products := v1.Group(("/products"))
		{
			products.POST("", productHandler.CreateProduct)
		}
	}

	return r
}
