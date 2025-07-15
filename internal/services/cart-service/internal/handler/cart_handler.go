package handler

import (
	"github.com/gin-gonic/gin"
	dependency_container "github.com/toji-dev/go-shop/internal/services/cart-service/internal/dependency-container"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/usecase"
)

type CartHandler struct {
	cartUseCase usecase.CartUseCase
}

func NewCartHandler(dependencyContainer *dependency_container.DependencyContainer) *CartHandler {
	return &CartHandler{
		cartUseCase: usecase.NewCartUseCase(),
	}
}

func (h *CartHandler) GetCart(c *gin.Context) {
}
