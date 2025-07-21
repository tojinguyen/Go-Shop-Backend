package router

import (
	"github.com/gin-gonic/gin"
	common_middleware "github.com/toji-dev/go-shop/internal/pkg/middleware"
	dependency_container "github.com/toji-dev/go-shop/internal/services/order-service/internal/dependency-container"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/handler"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/middleware"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func Init(router *gin.Engine, dependencyContainer *dependency_container.DependencyContainer) {
	cfg := dependencyContainer.GetConfig()

	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(common_middleware.ErrorHandler())

	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "order-service",
		})
	})

	orderHandler := handler.NewOrderHandler(dependencyContainer)

	v1 := router.Group("/api/v1")
	{
		orders := v1.Group("/orders")
		orders.Use(middleware.AuthHeaderMiddleware())
		{
			orders.GET("", orderHandler.GetOrdersByOwnerID)
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/:order_id", orderHandler.GetOrderByID)
		}
	}
}
