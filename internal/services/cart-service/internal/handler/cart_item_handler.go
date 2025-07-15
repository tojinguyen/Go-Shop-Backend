package handler

import (
	"github.com/gin-gonic/gin"
	dependency_container "github.com/toji-dev/go-shop/internal/services/cart-service/internal/dependency-container"
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
}

func (h *CartItemHandler) UpdateCartItem(c *gin.Context) {
}

func (h *CartItemHandler) RemoveCartItem(c *gin.Context) {
}