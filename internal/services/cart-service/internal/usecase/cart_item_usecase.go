package usecase

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/pkg/constant"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/dto"
	grpc "github.com/toji-dev/go-shop/internal/services/cart-service/internal/grpc/adapter"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/repository"
)

type cartItemUseCase struct {
	cartRepo       repository.CartRepository
	productAdapter grpc.ProductServiceAdapter
}

type CartItemUseCase interface {
	AddItemToCart(ctx *gin.Context, req dto.AddCartItemRequest) error
	UpdateCartItem()
	RemoveCartItem()
}

func NewCartItemUseCase(cartRepo repository.CartRepository, productAdapter grpc.ProductServiceAdapter) CartItemUseCase {
	return &cartItemUseCase{
		cartRepo:       cartRepo,
		productAdapter: productAdapter,
	}
}

func (uc *cartItemUseCase) AddItemToCart(ctx *gin.Context, req dto.AddCartItemRequest) error {
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return apperror.NewBadRequest("Invalid product ID format", err)
	}

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

	cart, err := uc.cartRepo.GetCartByOwnerID(ctx, userIDUUID)
	if err != nil {
		if appErr, ok := err.(*apperror.AppError); ok && appErr.Type == apperror.TypeNotFound {
			log.Printf("Cart not found for user %s. Creating a new one.", userID)
			cart = domain.NewCart(userIDUUID)
		} else {
			return apperror.NewInternal(fmt.Sprintf("Failed to get cart: %v", err))
		}
	}

	if err := cart.AddItem(productID, req.Quantity); err != nil {
		log.Printf("Failed to add item to cart: %v", err)
		return apperror.NewBadRequest("Failed to add item to cart", err)
	}

	if err := uc.cartRepo.Save(ctx, cart); err != nil {
		log.Printf("Failed to save cart: %v", err)
		return apperror.NewInternal(fmt.Sprintf("Failed to save cart: %v", err))
	}

	return nil
}

func (uc *cartItemUseCase) UpdateCartItem() {
}

func (uc *cartItemUseCase) RemoveCartItem() {
}
