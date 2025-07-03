package router

import (
	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/handlers"
)

// SetupRoutes sets up all the routes for the shop service
func SetupRoutes(r *gin.Engine) {
	// Initialize handlers
	shopHandler := handlers.NewShopHandler()

	// Health check
	r.GET("/health", shopHandler.HealthCheck)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Shop routes
		shops := v1.Group("/shops")
		{
			shops.GET("", shopHandler.GetShops)
			shops.POST("", shopHandler.CreateShop)
			shops.GET("/:id", shopHandler.GetShop)
			shops.PUT("/:id", shopHandler.UpdateShop)
			shops.DELETE("/:id", shopHandler.DeleteShop)
		}
	}
}
