package deleteshop

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/toji-dev/go-shop/internal/pkg/response"
)

// APIHandler handles HTTP requests for deleting shops
type APIHandler struct {
	commandHandler DeleteShopCommandHandler
	validator      *validator.Validate
}

// NewAPIHandler creates a new APIHandler for deleting shops
func NewAPIHandler(commandHandler DeleteShopCommandHandler) *APIHandler {
	return &APIHandler{
		commandHandler: commandHandler,
		validator:      validator.New(),
	}
}

// DeleteShop handles DELETE /api/v1/shops/{id}
func (h *APIHandler) DeleteShop(c *gin.Context) {
	shopID := c.Param("id")
	if shopID == "" {
		response.BadRequest(c, "SHOP_ID_REQUIRED", "Shop ID is required", "")
		return
	}

	command := DeleteShopCommand{
		ID: shopID,
	}

	// Validate the command
	if err := h.validator.Struct(command); err != nil {
		response.UnprocessableEntity(c, "VALIDATION_ERROR", "Validation failed", err.Error())
		return
	}

	// Handle the command
	err := h.commandHandler.Handle(c.Request.Context(), command)
	if err != nil {
		response.InternalServerError(c, "DELETE_SHOP_ERROR", "Failed to delete shop")
		return
	}

	response.Success(c, "Shop deleted successfully", nil)
}
