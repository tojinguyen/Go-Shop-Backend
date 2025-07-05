package api

import (
	"github.com/gin-gonic/gin"
	createshop "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/commands/create_shop"
	deleteshop "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/commands/delete_shop"
	updateshop "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/commands/update_shop"
	getshop "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/queries/get_shop"
	getshops "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/queries/get_shops"
)

// RegisterShopRoutes registers shop-related routes
func RegisterShopRoutes(
	r *gin.Engine,
	createShopAPIHandler *createshop.APIHandler,
	getShopAPIHandler *getshop.APIHandler,
	getShopsAPIHandler *getshops.APIHandler,
	updateShopAPIHandler *updateshop.APIHandler,
	deleteShopAPIHandler *deleteshop.APIHandler,
) {
	// Shop management routes
	shops := r.Group("/api/v1/shops")
	{
		// Create shop
		shops.POST("", createShopAPIHandler.CreateShop)

		// Get all shops for a user
		shops.GET("", getShopsAPIHandler.GetShops)

		// Get shop details by ID
		shops.GET("/:id", getShopAPIHandler.GetShop)

		// Update shop
		shops.PUT("/:id", updateShopAPIHandler.UpdateShop)

		// Delete shop
		shops.DELETE("/:id", deleteShopAPIHandler.DeleteShop)
	}
}
