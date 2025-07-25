package createshop

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/response"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/constant"
)

// APIHandler handles HTTP requests for create shop feature
type APIHandler struct {
	handler *Handler
}

// NewAPIHandler creates a new API handler for create shop
func NewAPIHandler(handler *Handler) *APIHandler {
	return &APIHandler{
		handler: handler,
	}
}

// CreateShop handles POST /shops
func (h *APIHandler) CreateShop(c *gin.Context) {
	log.Println("CreateShop API called")
	var req CreateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, constant.ErrorCodeValidation, "Invalid request data", err.Error())
		return
	}

	cmd := req.ToCommand()
	result, err := h.handler.Handle(c, cmd)

	if err != nil {
		log.Printf("Error creating shop: %v", err)
		response.InternalServerError(c, constant.ErrorCodeInternalError, "Failed to create shop")
		return
	}

	response.Created(c, constant.StatusCreated, result)
}
