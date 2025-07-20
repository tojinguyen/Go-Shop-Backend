package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/pkg/response"
	dependency_container "github.com/toji-dev/go-shop/internal/services/cart-service/internal/dependency-container"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/usecase"
)

type CartItemHandler struct {
	cartItemUseCase usecase.CartItemUseCase
}

func NewCartItemHandler(dependencyContainer *dependency_container.DependencyContainer) *CartItemHandler {
	return &CartItemHandler{
		cartItemUseCase: dependencyContainer.GetCartItemUseCase(),
	}
}

func (h *CartItemHandler) AddItemToCart(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, string(apperror.CodeUnauthorized), "User ID is required")
		return
	}

	var request dto.AddCartItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, string(apperror.CodeBadRequest), "Invalid request data", err.Error())
		return
	}
}

func (h *CartItemHandler) UpdateCartItem(c *gin.Context) {
}

func (h *CartItemHandler) RemoveCartItem(c *gin.Context) {
}
