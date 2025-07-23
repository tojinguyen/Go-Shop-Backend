package usecase

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/grpc/adapter"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/repository"
	product_v1 "github.com/toji-dev/go-shop/proto/gen/go/product/v1"
	shop_v1 "github.com/toji-dev/go-shop/proto/gen/go/shop/v1"
)

type OrderUsecase interface {
	CreateOrder(ctx *gin.Context, userId string, req dto.CreateOrderRequest) (*dto.OrderResponse, error)
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

func (u *orderUsecase) CreateOrder(ctx *gin.Context, userId string, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	order := &dto.OrderResponse{
		ShopID:            req.ShopID,
		ShippingAddressID: req.ShippingAddressID,
		BillingAddressID:  req.ShippingAddressID,
		PromotionID:       req.PromotionID,
		Note:              req.Note,
	}

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

	// Validate products in order items
	productIDs := make([]string, 0, len(req.Items))
	quantityMap := make(map[string]int32)
	for i, item := range req.Items {
		if item.ProductID == "" {
			return nil, apperror.NewBadRequest("Product ID cannot be empty", errors.New("product_id is required"))
		}

		if item.Quantity < 1 {
			return nil, apperror.NewBadRequest("Quantity must be at least 1", errors.New("quantity must be at least 1"))
		}
		productIDs[i] = item.ProductID
		quantityMap[item.ProductID] = int32(item.Quantity)
	}

	productInfoReponse, err := u.productServiceAdapter.GetProductsInfo(ctx, &product_v1.GetProductsInfoRequest{
		ProductIds: productIDs,
	})

	if !productInfoReponse.Valid || productInfoReponse.Products == nil || len(productInfoReponse.Products) != len(req.Items) {
		return nil, apperror.NewBadRequest("Invalid product information", errors.New("product information is invalid"))
	}

	totalAmount := 0.0
	for _, product := range productInfoReponse.Products {
		productPrice := product.Quantity * product.Price
		totalAmount += float64(productPrice)
	}

	// Validate Promotions
	if req.PromotionID != nil && *req.PromotionID != "" {
		promotionReq := &shop_v1.CalculatePromotionRequest{
			ShopId:        req.ShopID,
			UserId:        userId,
			PromotionCode: *req.PromotionID,
			TotalAmount:   int32(totalAmount),
		}

		promotionRes, err := u.shopServiceAdapter.CalculatePromotion(ctx, promotionReq)
		if err != nil {
			return nil, apperror.NewInternal(fmt.Sprintf("Failed to calculate promotion: %s", err.Error()))
		}

		if !promotionRes.Eligible {
			return nil, apperror.NewBadRequest("Invalid promotion code", errors.New("promotion code is not eligible"))
		}

		totalAmount -= float64(promotionRes.Discount)
		order.DiscountAmount = float64(promotionRes.Discount)
		order.TotalAmount = totalAmount
	}

	// Call to product service to reserve stock

	// Create order in repository

	return order, nil
}
