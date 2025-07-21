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
		log.Printf("Invalid product ID format: %v", err)
		return apperror.NewBadRequest("Invalid product ID format", err)
	}

	info, productInfoErr := uc.productAdapter.GetProductInfo(ctx, req.ProductID)
	if productInfoErr != nil {
		log.Printf("Failed to get product info: %v", productInfoErr)
		return apperror.NewInternal(fmt.Sprintf("failed to get product info: %v", productInfoErr))
	}

	if info.Exists == false {
		log.Printf("Product not found: %s", req.ProductID)
		return apperror.NewNotFound("Product", req.ProductID)
	}

	userID := ctx.GetString(constant.ContextKeyUserID)
	if userID == "" {
		log.Printf("User ID not found in context")
		return apperror.NewUnauthorized("User ID is required")
	}

	// Kiem tra xem gio hang da ton tai chua
	// Neu chua ton tai, tao moi gio hang
	userIDUUID, uuidErr := uuid.Parse(userID)
	if uuidErr != nil {
		log.Printf("Invalid user ID format: %v", uuidErr)
		return apperror.NewUnauthorized("Invalid user ID format")
	}

	cart, cartErr := uc.cartRepo.GetCartByOwnerID(ctx, userIDUUID)
	if cartErr != nil {
		if cartErr.Type == apperror.TypeNotFound {
			log.Printf("No cart found for user %s. Creating a new one.", userID)
			cart = domain.NewCart(userIDUUID)
		} else {
			log.Printf("Failed to get cart by owner ID %s: %v", userID, cartErr)
			return apperror.NewInternal(fmt.Sprintf("Failed to get cart by owner ID %s: %v", userID, cartErr))
		}
	}

	if addItemErr := cart.AddItem(productID, req.Quantity); addItemErr != nil {
		log.Printf("Failed to add item to cart: %v", addItemErr)
		return apperror.NewBadRequest("Failed to add item to cart", addItemErr)
	}

	if saveErr := uc.cartRepo.Save(ctx, cart); saveErr != nil {
		log.Printf("Failed to save cart: %v", saveErr)
		return apperror.NewInternal(fmt.Sprintf("Failed to save cart: %v", saveErr))
	}

	return nil
}
