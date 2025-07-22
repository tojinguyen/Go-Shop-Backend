package usecase

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/grpc/adapter"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/repository"
)

type OrderUsecase interface {
	CreateOrder(ctx *gin.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error)
}

type orderUsecase struct {
	orderRepo          repository.OrderRepository
	shopServiceAdapter adapter.ShopServiceAdapter
}

func NewOrderUsecase(orderRepo repository.OrderRepository, shopServiceAdapter adapter.ShopServiceAdapter) OrderUsecase {
	return &orderUsecase{orderRepo: orderRepo, shopServiceAdapter: shopServiceAdapter}
}

func (u *orderUsecase) CreateOrder(ctx *gin.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	shopID := req.ShopID

	isShopExists, err := u.shopServiceAdapter.CheckShopExists(ctx, shopID)
	if err != nil {
		return nil, apperror.NewInternal(fmt.Sprintf("Failed to check shop existence: %s", err.Error()))
	}

	if !isShopExists {
		return nil, apperror.NewNotFound("Shop", shopID)
	}

	// Check promotion eligibility
	if *req.PromotionID != "" {
		// Logic to check promotion code validity
		// For now, we assume the promotion code is valid
	}

	return nil, nil
}
