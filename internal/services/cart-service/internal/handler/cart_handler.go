package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	constant "github.com/toji-dev/go-shop/internal/pkg/constant"
	"github.com/toji-dev/go-shop/internal/pkg/response"
	dependency_container "github.com/toji-dev/go-shop/internal/services/cart-service/internal/dependency-container"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/usecase"
)

type CartHandler struct {
	cartUseCase usecase.CartUseCase
}

func NewCartHandler(dependencyContainer *dependency_container.DependencyContainer) *CartHandler {
	return &CartHandler{
		cartUseCase: dependencyContainer.GetCartUseCase(),
	}
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetString(constant.ContextKeyUserID)
	cart, err := h.cartUseCase.GetCart(c, userID)
	if err != nil {
		if apperror.GetType(err) == apperror.TypeNotFound {
			response.NotFound(c, "cart", userID)
			return
		}
		response.InternalServerError(c, "Failed to get cart", err.Error())
		return
	}
	response.Success(c, "success get cart", cart)
}

func (h *CartHandler) DeleteCart(c *gin.Context) {
	cartID := c.Param("id")
	if err := h.cartUseCase.DeleteCart(c, cartID); err != nil {
		if apperror.GetType(err) == apperror.TypeNotFound {
			response.NotFound(c, "cart", cartID)
			return
		}
		response.InternalServerError(c, "Failed to delete cart", err.Error())
		return
	}
	response.Success(c, "success delete cart", nil)
}
