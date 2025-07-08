package getpromotions

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

func (h *APIHandler) GetPromotions(c *gin.Context) {
	shopIDStr := c.Param("id")
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_SHOP_ID", "Invalid shop ID format", err.Error())
		return
	}

	promos, err := h.handler.Handle(c.Request.Context(), shopID)
	if err != nil {
		response.InternalServerError(c, "GET_PROMOTIONS_FAILED", err.Error())
		return
	}

	response.Success(c, "Promotions retrieved successfully", promos)
}
