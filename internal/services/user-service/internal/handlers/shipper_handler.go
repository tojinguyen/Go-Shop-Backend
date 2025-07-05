package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/response"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/services"
)

// ShipperHandler handles shipper-related requests
type ShipperHandler struct {
	shipperService *services.ShipperService
}

// NewShipperHandler creates a new shipper handler
func NewShipperHandler(sc container.ServiceContainer) *ShipperHandler {
	return &ShipperHandler{
		shipperService: services.NewShipperService(&sc),
	}
}

// RegisterShipper handles shipper registration
// @Summary Register as shipper
// @Description Register the authenticated user as a shipper
// @Tags shipper
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ShipperRegisterRequest true "Shipper registration request"
// @Success 201 {object} map[string]interface{} "Shipper registered successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /shipper/register [post]
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

// GetShipperProfile gets the current shipper's profile
// @Summary Get current shipper profile
// @Description Get the authenticated shipper's profile information
// @Tags shipper
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Shipper profile retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Shipper profile not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /shipper/profile [get]
func (h *ShipperHandler) GetShipperProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "UNAUTHORIZED", "User ID not found in context")
		return
	}

	shipper, err := h.shipperService.GetShipperProfile(c, userID.(string))
	if err != nil {
		response.InternalServerError(c, "INTERNAL_SERVER_ERROR", "Failed to retrieve shipper profile")
		return
	}

	if shipper == nil {
		response.NotFound(c, "NOT_FOUND", "Shipper profile not found")
		return
	}

	response.Success(c, "Shipper profile retrieved successfully", shipper)
}

// GetShipperProfileByID gets a shipper profile by ID
// @Summary Get shipper profile by ID
// @Description Get a shipper's profile by their user ID
// @Tags shipper
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{} "Shipper profile retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 404 {object} map[string]interface{} "Shipper profile not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /shipper/profile/{id} [get]
func (h *ShipperHandler) GetShipperProfileByID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "INVALID_REQUEST", "User ID is required", "User ID cannot be empty")
		return
	}

	shipper, err := h.shipperService.GetShipperProfile(c, userID)
	if err != nil {
		response.InternalServerError(c, "INTERNAL_SERVER_ERROR", "Failed to retrieve shipper profile")
		return
	}

	if shipper == nil {
		response.NotFound(c, "NOT_FOUND", "Shipper profile not found")
		return
	}

	response.Success(c, "Shipper profile retrieved successfully", shipper)
}

// UpdateShipperProfile updates the current shipper's profile
// @Summary Update shipper profile
// @Description Update the authenticated shipper's profile information
// @Tags shipper
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ShipperUpdateRequest true "Shipper update request"
// @Success 200 {object} map[string]interface{} "Shipper profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /shipper/profile [put]
func (h *ShipperHandler) UpdateShipperProfile(c *gin.Context) {
	var req dto.ShipperUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request payload", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "UNAUTHORIZED", "User ID not found in context")
		return
	}

	shipper, err := h.shipperService.UpdateShipperProfile(c, userID.(string), &req)
	if err != nil {
		response.InternalServerError(c, "INTERNAL_SERVER_ERROR", "Failed to update shipper profile")
		return
	}

	response.Success(c, "Shipper profile updated successfully", shipper)
}

// DeleteShipperProfile deletes the current shipper's profile
// @Summary Delete shipper profile
// @Description Delete the authenticated shipper's profile
// @Tags shipper
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Shipper profile deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /shipper/profile [delete]
func (h *ShipperHandler) DeleteShipperProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "UNAUTHORIZED", "User ID not found in context")
		return
	}

	err := h.shipperService.DeleteShipperProfile(c, userID.(string))
	if err != nil {
		response.InternalServerError(c, "INTERNAL_SERVER_ERROR", "Failed to delete shipper profile")
		return
	}

	response.Success(c, "Shipper profile deleted successfully", nil)
}
