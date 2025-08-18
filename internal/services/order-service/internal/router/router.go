package router

import (
	"github.com/gin-gonic/gin"
	common_middleware "github.com/toji-dev/go-shop/internal/pkg/middleware"
	dependency_container "github.com/toji-dev/go-shop/internal/services/order-service/internal/dependency-container"
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

	p := ginprometheus.NewPrometheus("go")
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.FullPath()
		if url == "" {
			url = "unknown"
		}
		return url
	}
	p.Use(router)

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(common_middleware.ErrorHandler())
	router.Use(common_middleware.OtelTracingMiddleware(cfg.App.Name))
	router.Use(common_middleware.AuthTokenMiddleware(dependencyContainer.GetJwtService()))

	orderHandler := dependencyContainer.GetOrderHandler()

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
