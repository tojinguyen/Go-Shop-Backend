package getshop

import (
	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/response"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/constant"
)

type APIHandler struct {
	handler *Handler
}

func NewAPIHandler(handler *Handler) *APIHandler {
	return &APIHandler{
		handler: handler,
	}
}

func (h *APIHandler) GetShop(c *gin.Context) {
	shopID := c.Param("id")
	if shopID == "" {
		response.BadRequest(c, constant.ErrorCodeValidation, "Shop ID is required", "shop_id is empty")
		return
	}

	resp, err := h.handler.Handle(c.Request.Context(), shopID)
	if err != nil {
		response.InternalServerError(c, constant.ErrorCodeInternalError, "Failed to get shop details")
		return
	}

	response.Success(c, constant.StatusOK, resp)
}
