package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/response"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/services"
)

type ShipperHandler struct {
	shipperService *services.ShipperService
}

func NewShipperHandler(sc container.ServiceContainer) *ShipperHandler {
	return &ShipperHandler{
		shipperService: services.NewShipperService(&sc),
	}
}

func (h *ShipperHandler) RegisterShipper(c *gin.Context) {
	var req dto.ShipperRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request payload", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "UNAUTHORIZED", "User ID not found in context")
		return
	}

	shipper, err := h.shipperService.RegisterShipper(c, userID.(string), &req)

	if err != nil {
		response.InternalServerError(c, "INTERNAL_SERVER_ERROR", "Failed to register shipper")
		return
	}

	response.Created(c, "Shipper registered successfully", shipper)
}

func (h *ShipperHandler) GetShipperProfile(c *gin.Context) {
	// Implementation pending
}

func (h *ShipperHandler) GetShipperProfileByID(c *gin.Context) {
	// Implementation pending
}

func (h *ShipperHandler) UpdateShipperProfile(c *gin.Context) {
	// Implementation pending
}

func (h *ShipperHandler) DeleteShipperProfile(c *gin.Context) {
	// Implementation pending
}
