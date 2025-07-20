package usecase

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/pkg/converter"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/repository"
)

type cartUseCase struct {
	repo repository.CartRepository
}

type CartUseCase interface {
	GetCart(ctx *gin.Context, userID string) (*domain.Cart, *apperror.AppError)
	DeleteCart(ctx *gin.Context, cartID string) *apperror.AppError
}

func NewCartUseCase(repo repository.CartRepository) CartUseCase {
	return &cartUseCase{repo: repo}
}

func (uc *cartUseCase) GetCart(ctx *gin.Context, userID string) (*domain.Cart, *apperror.AppError) {
	cart, err := uc.repo.GetCartByUserID(ctx, converter.StringToUUID(userID))
	if err != nil {
		if apperror.GetType(err) == apperror.TypeNotFound {
			return nil, apperror.NewNotFound("cart", userID)
		}
		return nil, apperror.NewInternal("Failed to get cart: " + fmt.Sprintf("%v", err))
	}
	return cart, nil
}

func (uc *cartUseCase) DeleteCart(ctx *gin.Context, cartID string) *apperror.AppError {
	if err := uc.repo.DeleteCart(ctx, converter.StringToUUID(cartID)); err != nil {
		if apperror.GetType(err) == apperror.TypeNotFound {
			return apperror.NewNotFound("cart", cartID)
		}
		return apperror.NewInternal("Failed to delete cart: " + fmt.Sprintf("%v", err))
	}
	return nil
}
