package api

import (
	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/handlers"
)

// RegisterShopRoutes registers shop-related routes
func RegisterShopRoutes(r *gin.Engine, shopHandler *handlers.ShopHandler) {
	// Shop management routes
	shops := r.Group("/api/v1/shops")
	{
		// Basic CRUD operations
		shops.POST("", shopHandler.CreateShop)        // Create shop
		shops.GET("", shopHandler.ListShops)          // List shops with pagination
		shops.GET("/search", shopHandler.SearchShops) // Search shops
		shops.GET("/:id", shopHandler.GetShop)        // Get shop by ID
		shops.PUT("/:id", shopHandler.UpdateShop)     // Update shop
		shops.DELETE("/:id", shopHandler.DeleteShop)  // Delete shop

		// Shop management operations
		shops.PUT("/:id/activate", shopHandler.ActivateShop) // Activate shop
		shops.PUT("/:id/ban", shopHandler.BanShop)           // Ban shop

		// Owner-specific routes
		shops.GET("/owner/:owner_id", shopHandler.GetShopsByOwner) // Get shops by owner
	}
}
