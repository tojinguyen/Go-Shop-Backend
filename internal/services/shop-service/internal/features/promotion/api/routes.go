package api

import (
	"github.com/gin-gonic/gin"
	createpromotion "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/promotion/commands/create_promotion"
	deletepromotion "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/promotion/commands/delete_promotion"
	updatepromotion "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/promotion/commands/update_promotion"
	getpromotionbyid "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/promotion/queries/get_promotion_by_id"
	getpromotions "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/promotion/queries/get_promotions"
)

func RegisterPromotionRoutes(
	r *gin.Engine,
	createPromoAPIHandler *createpromotion.APIHandler,
	getPromosAPIHandler *getpromotions.APIHandler,
	getPromoByIDAPIHandler *getpromotionbyid.APIHandler,
	updatePromoAPIHandler *updatepromotion.APIHandler,
	deletePromoAPIHandler *deletepromotion.APIHandler,
) {
	// Shop management routes
	promotions := r.Group("/api/v1/shops/{id}/promotions")
	{
		promotions.POST("", createPromoAPIHandler.CreatePromotion)
		promotions.GET("", getPromosAPIHandler.GetPromotions)
		promotions.GET("/:promo_id", getPromoByIDAPIHandler.GetPromotionByID)
		promotions.PUT("/:promo_id", updatePromoAPIHandler.UpdatePromotion)
		promotions.DELETE("/:promo_id", deletePromoAPIHandler.DeletePromotion)
	}
}
