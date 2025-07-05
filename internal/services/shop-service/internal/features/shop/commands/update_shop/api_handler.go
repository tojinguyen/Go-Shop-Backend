package updateshop

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/toji-dev/go-shop/internal/pkg/response"
)

// APIHandler handles HTTP requests for updating shops
type APIHandler struct {
	commandHandler UpdateShopCommandHandler
	validator      *validator.Validate
}

// NewAPIHandler creates a new APIHandler for updating shops
func NewAPIHandler(commandHandler UpdateShopCommandHandler) *APIHandler {
	return &APIHandler{
		commandHandler: commandHandler,
		validator:      validator.New(),
	}
}

// UpdateShop handles PUT /api/v1/shops/{id}
func (h *APIHandler) UpdateShop(c *gin.Context) {
	shopID := c.Param("id")
	if shopID == "" {
		response.BadRequest(c, "SHOP_ID_REQUIRED", "Shop ID is required", "")
		return
	}

	var command UpdateShopCommand
	if err := c.ShouldBindJSON(&command); err != nil {
		response.BadRequest(c, "INVALID_JSON", "Invalid JSON format", err.Error())
		return
	}

	// Set the ID from URL parameter
	command.ID = shopID

	// Validate the command
	if err := h.validator.Struct(command); err != nil {
		response.UnprocessableEntity(c, "VALIDATION_ERROR", "Validation failed", err.Error())
		return
	}

	// Handle the command
	updatedShop, err := h.commandHandler.Handle(c.Request.Context(), command)
	if err != nil {
		response.InternalServerError(c, "UPDATE_SHOP_ERROR", "Failed to update shop")
		return
	}

	response.Success(c, "Shop updated successfully", updatedShop)
}
