package router

import (
	"github.com/gin-gonic/gin"
	dependency_container "github.com/toji-dev/go-shop/internal/services/cart-service/internal/dependency-container"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/handler"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/middleware"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func SetupRoutes(r *gin.Engine, dependencyContainer *dependency_container.DependencyContainer) {
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
	p.Use(r)

	cartHandler := handler.NewCartHandler(dependencyContainer)
	cartItemHandler := handler.NewCartItemHandler(dependencyContainer)

	v1 := r.Group("/api/v1")
	{
		cart := v1.Group("/carts")
		cart.Use(middleware.AuthHeaderMiddleware())
		{
			cart.GET("", cartHandler.GetCart)
			cart.DELETE("", cartHandler.DeleteCartByOwnerID)
			cart.POST("/items", cartItemHandler.UpdateItemsInCart)
		}
	}
}
