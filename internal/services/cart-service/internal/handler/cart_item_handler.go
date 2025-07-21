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

func (h *CartItemHandler) UpdateItemsInCart(c *gin.Context) {
	var request dto.AddCartItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, string(apperror.CodeBadRequest), "Invalid request data", err.Error())
		return
	}

	if request.ProductID == "" {
		response.BadRequest(c, string(apperror.CodeBadRequest), "Cart ID is required", "")
		return
	}

	if request.Quantity <= 0 {
		response.BadRequest(c, string(apperror.CodeBadRequest), "Quantity must be greater than zero", "")
		return
	}

	err := h.cartItemUseCase.AddItemToCart(c, request)

	if err != nil {
		if apperror.GetType(err) == apperror.TypeNotFound {
			response.NotFound(c, string(apperror.CodeNotFound), "Product not found")
		} else {
			response.InternalServerError(c, string(apperror.CodeInternal), "Failed to add item to cart")
		}
		return
	}

	response.Success(c, "success add item to cart", nil)
}
