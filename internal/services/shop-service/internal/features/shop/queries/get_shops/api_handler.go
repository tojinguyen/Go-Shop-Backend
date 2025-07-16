package getshops

import (
	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/constant"
	"github.com/toji-dev/go-shop/internal/pkg/response"
)

// APIHandler handles HTTP requests for getting shops
type APIHandler struct {
	queryHandler *QueryHandler
}

// NewAPIHandler creates a new APIHandler for getting shops
func NewAPIHandler(queryHandler *QueryHandler) *APIHandler {
	return &APIHandler{
		queryHandler: queryHandler,
	}
}

// GetShops handles GET /api/v1/shops
func (h *APIHandler) GetShops(c *gin.Context) {
	// Get owner_id from query parameter or JWT token
	ownerID := c.Query("owner_id")
	if ownerID == "" {
		// If not provided in query, try to get from authenticated user context
		if userID, exists := c.Get(constant.ContextKeyUserID); exists {
			ownerID = userID.(string)
		} else {
			response.BadRequest(c, "OWNER_ID_REQUIRED", "Owner ID is required", "")
			return
		}
	}

	query := GetShopsQuery{
		OwnerID: ownerID,
	}

	shops, err := h.queryHandler.Handle(c.Request.Context(), query)
	if err != nil {
		response.InternalServerError(c, "GET_SHOPS_ERROR", "Failed to get shops")
		return
	}

	response.Success(c, "Shops retrieved successfully", shops)
}
