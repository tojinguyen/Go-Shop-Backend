package api

import (
	"github.com/gin-gonic/gin"
	createshop "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/commands/create_shop"
	getshop "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/queries/get_shop"
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

// RegisterGetShopRoutes registers routes for getting shop details
func RegisterGetShopRoutes(r *gin.Engine, getShopAPIHandler *getshop.APIHandler) {
	// Shop retrieval routes
	shops := r.Group("/api/v1/shops")
	{
		// Get shop details by ID
		shops.GET("/:id", getShopAPIHandler.GetShop) // Get shop by ID
	}
}
