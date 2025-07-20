package usecase

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/pkg/constant"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/dto"
	grpc "github.com/toji-dev/go-shop/internal/services/cart-service/internal/grpc/adapter"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/repository"
)

type cartItemUseCase struct {
	cartItemRepo   repository.CartItemRepository
	cartRepo       repository.CartRepository
	productAdapter grpc.ProductServiceAdapter
}

type CartItemUseCase interface {
	AddItemToCart(ctx *gin.Context, req dto.AddCartItemRequest) error
	UpdateCartItem()
	RemoveCartItem()
}

func NewCartItemUseCase(cartItemRepo repository.CartItemRepository, cartRepo repository.CartRepository, productAdapter grpc.ProductServiceAdapter) CartItemUseCase {
	return &cartItemUseCase{
		cartItemRepo:   cartItemRepo,
		cartRepo:       cartRepo,
		productAdapter: productAdapter,
	}
}

func (uc *cartItemUseCase) AddItemToCart(ctx *gin.Context, req dto.AddCartItemRequest) error {
	info, err := uc.productAdapter.GetProductInfo(ctx, req.ProductID)
	if err != nil {
		return apperror.NewInternal(fmt.Sprintf("failed to get product info: %v", err))
	}

	if info.Exists == false {
		return apperror.NewNotFound("Product", req.ProductID)
	}

	userID := ctx.GetString(constant.ContextKeyUserID)
	if userID == "" {
		log.Printf("User ID not found in context")
		return apperror.NewUnauthorized("User ID is required")
	}

	// Kiem tra xem gio hang da ton tai chua
	// Neu chua ton tai, tao moi gio hang
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		return apperror.NewUnauthorized("Invalid user ID format")
	}

	cart, err := uc.cartRepo.GetCartByUserID(ctx, userIDUUID)
	if err != nil {
		if apperror.GetType(err) == apperror.TypeNotFound {
			return apperror.NewNotFound("Cart", userID)
		}
		return apperror.NewInternal(fmt.Sprintf("failed to get cart: %v", err))
	}

	// Kiem tra xem san pham da ton tai trong gio hang chua
	// Neu ton tai, cap nhat so luong
	productIDUUID, err := uuid.Parse(req.ProductID)
	if err != nil {
		log.Printf("Invalid product ID format: %v", err)
		return apperror.NewBadRequest("Invalid product ID format", err)
	}

	for _, item := range cart.Items {
		if item.ProductID == productIDUUID {
			item.Quantity += req.Quantity
			if err := uc.cartItemRepo.UpdateCartItem(ctx, &item); err != nil {
				return apperror.NewInternal(fmt.Sprintf("failed to update cart item: %v", err))
			}
			return nil
		}
	}

	// Neu khong ton tai, them san pham vao gio hang

	return nil
}

func (uc *cartItemUseCase) UpdateCartItem() {
}

func (uc *cartItemUseCase) RemoveCartItem() {
}
