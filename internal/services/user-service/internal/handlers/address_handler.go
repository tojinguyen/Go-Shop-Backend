package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/response"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/services"
)

// AddressHandler handles user address-related requests
type AddressHandler struct {
	addressService *services.AddressService
}

// NewAddressHandler creates a new address handler
func NewAddressHandler(sc container.ServiceContainer) *AddressHandler {
	return &AddressHandler{
		addressService: services.NewAddressService(&sc),
	}
}

// GetAddresses handles the request to get all addresses for a user
// @Summary Get user addresses
// @Description Get all addresses for the authenticated user
// @Tags address
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Addresses retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /addresses [get]
func (h *AddressHandler) GetAddresses(c *gin.Context) {
	var userID string
	if userIDParam, exists := c.Get("user_id"); exists {
		userID = userIDParam.(string)
	} else {
		response.Unauthorized(c, "UNAUTHORIZED", "User ID not found in context")
		return
	}

	addresses, err := h.addressService.GetAddressesByUserID(c, userID)
	if err != nil {
		response.InternalServerError(c, "GET_ADDRESSES_FAILED", "Failed to retrieve addresses")
		return
	}

	response.Success(c, "Addresses retrieved successfully", addresses)
}

// GetAddressByID handles the request to get a single address by ID
// @Summary Get address by ID
// @Description Get a specific address by its ID
// @Tags address
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Address ID"
// @Success 200 {object} map[string]interface{} "Address retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 404 {object} map[string]interface{} "Address not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /addresses/{id} [get]
func (h *AddressHandler) GetAddressByID(c *gin.Context) {
	// Implementation pending
	addressID := c.Param("id")
	if addressID == "" {
		log.Println("Address ID is required")
		response.BadRequest(c, "INVALID_REQUEST", "Address ID is required", "Address ID should not be empty")
		return
	}

	address, err := h.addressService.GetAddressByID(c, addressID)

	if err != nil {
		if err.Error() == "address not found" {
			log.Printf("Address with ID %s not found: %v", addressID, err)
			response.NotFound(c, "ADDRESS_NOT_FOUND", "Address with this ID does not exist")
			return
		}
		log.Printf("Failed to retrieve address by ID %s: %v", addressID, err)
		response.InternalServerError(c, "GET_ADDRESS_FAILED", "Failed to retrieve address")
		return
	}
	response.Success(c, "Address retrieved successfully", address)
}

// AddAddress handles the request to add a new address
// @Summary Create new address
// @Description Create a new address for the authenticated user
// @Tags address
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateAddressRequest true "Create address request"
// @Success 201 {object} map[string]interface{} "Address created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 409 {object} map[string]interface{} "Address already exists"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /addresses [post]
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

// UpdateAddress handles the request to update an existing address
// @Summary Update address
// @Description Update an existing address for the authenticated user
// @Tags address
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Address ID"
// @Param request body dto.UpdateAddressRequest true "Update address request"
// @Success 200 {object} map[string]interface{} "Address updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Address not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /addresses/{id} [put]
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	// Implementation pending
	var req dto.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request payload", err.Error())
		return
	}

	addressID := c.Param("id")
	if addressID == "" {
		response.BadRequest(c, "INVALID_REQUEST", "Address ID is required", "Address ID should not be empty")
		return
	}

	var userID string
	if userIDParam, exists := c.Get("user_id"); exists {
		userID = userIDParam.(string)
	} else {
		response.Unauthorized(c, "UNAUTHORIZED", "User ID not found in context")
		return
	}

	address, err := h.addressService.UpdateAddress(c, userID, addressID, req)
	if err != nil {
		if err.Error() == "address not found" {
			response.NotFound(c, "ADDRESS_NOT_FOUND", "Address with this ID does not exist")
			return
		}
		response.InternalServerError(c, "UPDATE_ADDRESS_FAILED", "Failed to update address")
		return
	}

	response.Success(c, "Address updated successfully", address)
}

// DeleteAddress handles the request to delete an address
// @Summary Delete address
// @Description Delete an address for the authenticated user
// @Tags address
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Address ID"
// @Success 200 {object} map[string]interface{} "Address deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Address not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /addresses/{id} [delete]
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	addressID := c.Param("id")
	if addressID == "" {
		response.BadRequest(c, "INVALID_REQUEST", "Address ID is required", "Address ID should not be empty")
		return
	}

	var userID string
	if userIDParam, exists := c.Get("user_id"); exists {
		userID = userIDParam.(string)
	} else {
		response.Unauthorized(c, "UNAUTHORIZED", "User ID not found in context")
		return
	}

	err := h.addressService.DeleteAddress(c, userID, addressID)
	if err != nil {
		if err.Error() == "address not found" {
			response.NotFound(c, "ADDRESS_NOT_FOUND", "Address with this ID does not exist")
			return
		}
		response.InternalServerError(c, "DELETE_ADDRESS_FAILED", "Failed to delete address")
		return
	}

	response.Success(c, "Address deleted successfully", nil)
}

// SetDefaultAddress handles the request to set an address as the default
// @Summary Set default address
// @Description Set an address as the default address for the authenticated user
// @Tags address
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Address ID"
// @Success 200 {object} map[string]interface{} "Default address set successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Address not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /addresses/{id}/default [put]
func (h *AddressHandler) SetDefaultAddress(c *gin.Context) {
	addressID := c.Param("id")
	if addressID == "" {
		response.BadRequest(c, "INVALID_REQUEST", "Address ID is required", "Address ID should not be empty")
		return
	}

	var userID string
	if userIDParam, exists := c.Get("user_id"); exists {
		userID = userIDParam.(string)
	} else {
		response.Unauthorized(c, "UNAUTHORIZED", "User ID not found in context")
		return
	}

	address, err := h.addressService.SetDefaultAddress(c, userID, addressID)
	if err != nil {
		if err.Error() == "address not found" {
			response.NotFound(c, "ADDRESS_NOT_FOUND", "Address with this ID does not exist")
			return
		}
		response.InternalServerError(c, "SET_DEFAULT_ADDRESS_FAILED", "Failed to set default address")
		return
	}

	response.Success(c, "Default address set successfully", address)
}
