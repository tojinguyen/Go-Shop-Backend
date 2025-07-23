package usecase

import (
	"errors"
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
	orderRepo             repository.OrderRepository
	shopServiceAdapter    adapter.ShopServiceAdapter
	productServiceAdapter adapter.ProductServiceAdapter
	userAdapter           adapter.UserServiceAdapter
}

func NewOrderUsecase(orderRepo repository.OrderRepository, shopServiceAdapter adapter.ShopServiceAdapter, productServiceAdapter adapter.ProductServiceAdapter, userAdapter adapter.UserServiceAdapter) OrderUsecase {
	return &orderUsecase{orderRepo: orderRepo, shopServiceAdapter: shopServiceAdapter, productServiceAdapter: productServiceAdapter, userAdapter: userAdapter}
}

func (u *orderUsecase) CreateOrder(ctx *gin.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {

	// Validate shop existence
	isShopExists, err := u.shopServiceAdapter.CheckShopExists(ctx, req.ShopID)
	if err != nil {
		return nil, apperror.NewInternal(fmt.Sprintf("Failed to check shop existence: %s", err.Error()))
	}

	if !isShopExists {
		return nil, apperror.NewNotFound("Shop", req.ShopID)
	}

	// Validate shipping address
	if req.ShippingAddressID == "" {
		return nil, apperror.NewBadRequest("Address cannot be empty", errors.New("shipping_address_id is required"))
	}

	address, err := u.userAdapter.GetAddressById(ctx, req.ShippingAddressID)
	if err != nil {
		return nil, apperror.NewInternal(fmt.Sprintf("Failed to get address: %s", err.Error()))
	}

	if address == nil || address.DeletedAt != nil {
		return nil, apperror.NewNotFound("Shipping address", req.ShippingAddressID)
	}

	return nil, nil
}
