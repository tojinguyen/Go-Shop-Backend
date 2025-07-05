package api

import (
	"github.com/gin-gonic/gin"
	createshop "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/commands/create_shop"
)

// RegisterShopRoutes registers shop-related routes
func RegisterShopRoutes(r *gin.Engine, createShopAPIHandler *createshop.APIHandler) {
	// Shop management routes
	shops := r.Group("/api/v1/shops")
	{
		// Create shop feature
		shops.POST("", createShopAPIHandler.CreateShop) // Create shop
	}
}
