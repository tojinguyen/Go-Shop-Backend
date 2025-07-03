package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
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
	// Implementation pending
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
