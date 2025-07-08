package createpromotion

import (
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

func (h *APIHandler) CreatePromotion(c *gin.Context) {
	shopIDStr := c.Param("id")
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_SHOP_ID", "Invalid shop ID format", err.Error())
		return
	}

	var req CreatePromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	promo, err := h.handler.Handle(c.Request.Context(), shopID, req)
	if err != nil {
		response.InternalServerError(c, "CREATE_PROMOTION_FAILED", err.Error())
		return
	}

	response.Created(c, "Promotion created successfully", promo)
}
