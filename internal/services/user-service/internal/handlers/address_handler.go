package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/response"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/services"
)

type AddressHandler struct {
	addressService *services.AddressService
}

func NewAddressHandler(sc container.ServiceContainer) *AddressHandler {
	return &AddressHandler{
		addressService: services.NewAddressService(&sc),
	}
}

// GetAddress handles the request to get all addresses for a user.
func (h *AddressHandler) GetAddress(c *gin.Context) {
	// Implementation pending
	addressID := c.Param("id")
	if addressID != "" {
		response.BadRequest(c, "INVALID_REQUEST", "Address ID is required", "Address ID should not be empty")
		return
	}

	address, err := h.addressService.GetAddressByID(c, addressID)

	if err != nil {
		if err.Error() == "address not found" {
			response.NotFound(c, "ADDRESS_NOT_FOUND", "Address with this ID does not exist")
			return
		}
		response.InternalServerError(c, "GET_ADDRESS_FAILED", "Failed to retrieve address")
		return
	}
	response.Success(c, "Address retrieved successfully", address)
}

// GetAddressByID handles the request to get a single address by ID.
func (h *AddressHandler) GetAddressByID(c *gin.Context) {
	// Implementation pending
}

// AddAddress handles the request to add a new address.
func (h *AddressHandler) AddAddress(c *gin.Context) {
	// Implementation pending
	var req dto.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request payload", err.Error())
		return
	}

	var userID string
	if userIDParam, exists := c.Get("user_id"); exists {
		userID = userIDParam.(string)
	} else {
		response.Unauthorized(c, "UNAUTHORIZED", "User ID not found in context")
		return
	}

	address, err := h.addressService.CreateAddress(c, userID, req)
	if err != nil {
		if err.Error() == "address already exists" {
			response.Conflict(c, "ADDRESS_ALREADY_EXISTS", "Address with this details already exists")
			return
		}
		response.InternalServerError(c, "CREATE_ADDRESS_FAILED", "Failed to create address")
		return
	}

	response.Created(c, "Address created successfully", address)
}

// UpdateAddress handles the request to update an existing address.
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	// Implementation pending
}

// DeleteAddress handles the request to delete an address.
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	// Implementation pending
}

// SetDefaultAddress handles the request to set an address as the default.
func (h *AddressHandler) SetDefaultAddress(c *gin.Context) {
	// Implementation pending
}
