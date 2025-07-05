package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/response"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/constant"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/service"
)

// ShopHandler handles HTTP requests for shop operations
type ShopHandler struct {
	shopService service.ShopService
}

// NewShopHandler creates a new shop handler
func NewShopHandler(shopService service.ShopService) *ShopHandler {
	return &ShopHandler{
		shopService: shopService,
	}
}

// CreateShop handles POST /shops
func (h *ShopHandler) CreateShop(c *gin.Context) {
	var req dto.CreateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, constant.ErrorCodeValidation, "Invalid request data", err.Error())
		return
	}

	cmd := req.ToCommand()
	result, err := h.shopService.CreateShop(c.Request.Context(), cmd)
	if err != nil {
		if err.Error() == "owner already has a shop" || err.Error() == "email already in use" {
			response.Conflict(c, constant.ErrorCodeDuplicate, err.Error())
			return
		}
		response.InternalServerError(c, constant.ErrorCodeInternalError, "Failed to create shop")
		return
	}

	response.Created(c, constant.StatusCreated, result)
}

// GetShop handles GET /shops/:id
func (h *ShopHandler) GetShop(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, constant.ErrorCodeValidation, "Shop ID is required", "")
		return
	}

	result, err := h.shopService.GetShopByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "shop not found" {
			response.NotFound(c, constant.ErrorCodeNotFound, "Shop not found")
			return
		}
		response.InternalServerError(c, constant.ErrorCodeInternalError, "Failed to get shop")
		return
	}

	response.Success(c, constant.StatusOK, result)
}

// GetShopsByOwner handles GET /shops/owner/:owner_id
func (h *ShopHandler) GetShopsByOwner(c *gin.Context) {
	ownerID := c.Param("owner_id")
	if ownerID == "" {
		response.BadRequest(c, constant.ErrorCodeValidation, "Owner ID is required", "")
		return
	}

	result, err := h.shopService.GetShopsByOwnerID(c.Request.Context(), ownerID)
	if err != nil {
		response.InternalServerError(c, constant.ErrorCodeInternalError, "Failed to get shops")
		return
	}

	response.Success(c, constant.StatusOK, result)
}

// UpdateShop handles PUT /shops/:id
func (h *ShopHandler) UpdateShop(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, constant.ErrorCodeValidation, "Shop ID is required", "")
		return
	}

	var req dto.UpdateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, constant.ErrorCodeValidation, "Invalid request data", err.Error())
		return
	}

	result, err := h.shopService.UpdateShop(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "shop not found" {
			response.NotFound(c, constant.ErrorCodeNotFound, "Shop not found")
			return
		}
		if err.Error() == "email already in use" {
			response.Conflict(c, constant.ErrorCodeDuplicate, err.Error())
			return
		}
		response.InternalServerError(c, constant.ErrorCodeInternalError, "Failed to update shop")
		return
	}

	response.Success(c, constant.StatusUpdated, result)
}

// DeleteShop handles DELETE /shops/:id
func (h *ShopHandler) DeleteShop(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, constant.ErrorCodeValidation, "Shop ID is required", "")
		return
	}

	err := h.shopService.DeleteShop(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "shop not found" {
			response.NotFound(c, constant.ErrorCodeNotFound, "Shop not found")
			return
		}
		response.InternalServerError(c, constant.ErrorCodeInternalError, "Failed to delete shop")
		return
	}

	response.Success(c, constant.StatusDeleted, nil)
}

// ListShops handles GET /shops
func (h *ShopHandler) ListShops(c *gin.Context) {
	// Parse pagination parameters
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	result, total, err := h.shopService.ListShops(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, constant.ErrorCodeInternalError, "Failed to list shops")
		return
	}

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	meta := &response.MetaInfo{
		Page:       page,
		PerPage:    limit,
		Total:      total,
		TotalPages: totalPages,
	}

	response.SuccessWithMeta(c, constant.StatusOK, result, meta)
}

// SearchShops handles GET /shops/search
func (h *ShopHandler) SearchShops(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, constant.ErrorCodeValidation, "Search query is required", "")
		return
	}

	// Parse pagination parameters
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	result, total, err := h.shopService.SearchShops(c.Request.Context(), query, limit, offset)
	if err != nil {
		response.InternalServerError(c, constant.ErrorCodeInternalError, "Failed to search shops")
		return
	}

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	meta := &response.MetaInfo{
		Page:       page,
		PerPage:    limit,
		Total:      total,
		TotalPages: totalPages,
	}

	response.SuccessWithMeta(c, constant.StatusOK, result, meta)
}

// ActivateShop handles PUT /shops/:id/activate
func (h *ShopHandler) ActivateShop(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, constant.ErrorCodeValidation, "Shop ID is required", "")
		return
	}

	err := h.shopService.ActivateShop(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "shop not found" {
			response.NotFound(c, constant.ErrorCodeNotFound, "Shop not found")
			return
		}
		response.InternalServerError(c, constant.ErrorCodeInternalError, "Failed to activate shop")
		return
	}

	response.Success(c, "Shop activated successfully", nil)
}

// BanShop handles PUT /shops/:id/ban
func (h *ShopHandler) BanShop(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, constant.ErrorCodeValidation, "Shop ID is required", "")
		return
	}

	err := h.shopService.BanShop(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "shop not found" {
			response.NotFound(c, constant.ErrorCodeNotFound, "Shop not found")
			return
		}
		response.InternalServerError(c, constant.ErrorCodeInternalError, "Failed to ban shop")
		return
	}

	response.Success(c, "Shop banned successfully", nil)
}
