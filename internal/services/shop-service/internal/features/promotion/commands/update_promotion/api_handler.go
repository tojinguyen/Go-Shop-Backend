package updatepromotion

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/pkg/response"
)

type APIHandler struct {
	handler *Handler
}

func NewAPIHandler(handler *Handler) *APIHandler {
	return &APIHandler{handler: handler}
}

func (h *APIHandler) UpdatePromotion(c *gin.Context) {
	shopIDStr := c.Param("id")
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_SHOP_ID", "Invalid shop ID format", err.Error())
		return
	}

	promoIDStr := c.Param("promo_id")
	promoID, err := uuid.Parse(promoIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_PROMOTION_ID", "Invalid promotion ID format", err.Error())
		return
	}

	var req UpdatePromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	err = h.handler.Handle(c.Request.Context(), promoID, shopID, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "PROMOTION_NOT_FOUND", fmt.Sprintf("Promotion with ID %s not found", promoIDStr))
			return
		}
		if strings.Contains(err.Error(), "does not belong") {
			response.Forbidden(c, "FORBIDDEN", "You do not have permission to update this promotion.")
			return
		}
		response.InternalServerError(c, "UPDATE_PROMOTION_FAILED", err.Error())
		return
	}

	response.Success(c, "Promotion updated successfully", nil)
}
