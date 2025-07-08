package getpromotionbyid

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

func (h *APIHandler) GetPromotionByID(c *gin.Context) {
	promoIDStr := c.Param("promo_id")
	promoID, err := uuid.Parse(promoIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_PROMOTION_ID", "Invalid promotion ID format", err.Error())
		return
	}

	promo, err := h.handler.Handle(c.Request.Context(), promoID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "PROMOTION_NOT_FOUND", fmt.Sprintf("Promotion with ID %s not found", promoIDStr))
			return
		}
		response.InternalServerError(c, "GET_PROMOTION_FAILED", err.Error())
		return
	}

	response.Success(c, "Promotion retrieved successfully", promo)
}
