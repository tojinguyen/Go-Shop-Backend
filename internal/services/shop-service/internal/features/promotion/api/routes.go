package api

import (
	"github.com/gin-gonic/gin"
	createpromotion "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/promotion/commands/create_promotion"
)

func RegisterPromotionRoutes(
	r *gin.Engine,
	createPromoAPIHandler *createpromotion.APIHandler,
) {
	// Shop management routes
	promotions := r.Group("/api/v1/shops/{id}/promotions")
	{
		promotions.POST("", createPromoAPIHandler.CreatePromotion)
	}
}
