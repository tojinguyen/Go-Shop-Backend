package usecase

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/grpc/adapter"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/repository"
	product_v1 "github.com/toji-dev/go-shop/proto/gen/go/product/v1"
	shop_v1 "github.com/toji-dev/go-shop/proto/gen/go/shop/v1"
)

type OrderUsecase interface {
	CreateOrder(ctx *gin.Context, userId string, req dto.CreateOrderRequest) (*domain.Order, error)
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

func (u *orderUsecase) CreateOrder(ctx *gin.Context, userId string, req dto.CreateOrderRequest) (*domain.Order, error) {
	// --- STAGE 1: VALIDATION ---
	// Validate shop, address, and product info before creating anything
	if err := u.validatePrerequisites(ctx, &req); err != nil {
		return nil, err
	}

	// Validate products in order items
	productIDs := make([]string, 0, len(req.Items))
	quantityMap := make(map[string]int32)
	for i, item := range req.Items {
		productIDs[i] = item.ProductID
		quantityMap[item.ProductID] = int32(item.Quantity)
	}

	productsInfo, err := u.productServiceAdapter.GetProductsInfo(ctx, &product_v1.GetProductsInfoRequest{
		ProductIds: productIDs,
	})

	if err != nil || !productsInfo.Valid || productsInfo.Products == nil || len(productsInfo.Products) != len(req.Items) {
		return nil, apperror.NewBadRequest("One or more products are invalid or unavailable", err)
	}

	// --- STAGE 2: CALCULATION & ORDER CREATION (PENDING) ---
	orderID := uuid.New().String()
	totalAmount := 0.0
	orderItems := make([]domain.OrderItem, len(productsInfo.Products))

	for i, p := range productsInfo.Products {
		price := float64(p.Price)
		totalAmount += price * float64(quantityMap[p.Id])
		orderItems[i] = domain.OrderItem{
			ProductID: p.Id,
			Quantity:  int(quantityMap[p.Id]),
			Price:     price,
		}
	}
	// Validate Promotions
	discountAmount := 0.0
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
		discountAmount = float64(promotionRes.Discount)
	}

	// Create the order with PENDING status FIRST
	pendingOrder := &domain.Order{
		ID:                orderID,
		OwnerID:           userId,
		ShopID:            req.ShopID,
		ShippingAddressID: req.ShippingAddressID,
		PromotionCode:     req.PromotionID,
		DiscountAmount:    discountAmount,
		TotalAmount:       totalAmount,
		Status:            domain.OrderStatusPENDING,
		Items:             orderItems,
	}

	createdOrder, err := u.orderRepo.CreateOrder(ctx, pendingOrder)

	if err != nil {
		return nil, apperror.NewInternal(fmt.Sprintf("Failed to create initial order record: %s", err.Error()))
	}

	// Call to product service to reserve stock
	var reserveItems []*product_v1.ReserveProduct
	for _, item := range req.Items {
		reserveItems = append(reserveItems, &product_v1.ReserveProduct{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
		})
	}

	reserveReq := &product_v1.ReserveProductsRequest{
		OrderId:  orderID,
		ShopId:   req.ShopID,
		Products: reserveItems,
	}

	reserveResp, err := u.productServiceAdapter.ReserveProducts(ctx, reserveReq)
	if err != nil {
		log.Printf("CRITICAL: ReserveProducts call failed for order %s. Marking as FAILED. Error: %v", orderID, err)
		u.orderRepo.UpdateOrderStatus(ctx, orderID, sqlc.OrderStatusFAILED)
		return nil, apperror.NewDependencyFailure(fmt.Sprintf("Failed to reserve products: %s", err.Error()))
	}

	// --- STAGE 4: FINALIZE ORDER (SAGA - COMMIT/ROLLBACK) ---
	if !reserveResp.Success {
		log.Printf("Failed to reserve products for order %s. Marking as FAILED.", orderID)
		u.orderRepo.UpdateOrderStatus(ctx, orderID, sqlc.OrderStatusFAILED)

		var errorDetails []string
		for _, status := range reserveResp.ProductStatuses {
			if !status.Success {
				errorDetails = append(errorDetails, fmt.Sprintf("product %s: %s", status.ProductId, status.Message))
			}
		}

		log.Printf("Failed to reserve products: %s", strings.Join(errorDetails, ", "))
		return nil, apperror.NewConflict("Failed to reserve all products", strings.Join(errorDetails, ", "))
	}

	finalOrder, err := u.orderRepo.UpdateOrderStatus(ctx, orderID, sqlc.OrderStatusPENDINGPAYMENT)

	if err != nil {
		log.Printf("CRITICAL: SAGA failure. Product stock reserved for order %s but failed to update order status. Manual intervention required. Error: %v", orderID, err)
		return nil, apperror.NewInternal("Failed to finalize order status after successful reservation.")
	}
	finalOrder.Items = createdOrder.Items

	return finalOrder, nil
}

func (u *orderUsecase) validatePrerequisites(ctx *gin.Context, req *dto.CreateOrderRequest) error {
	// Validate shop existence
	isShopExists, err := u.shopServiceAdapter.CheckShopExists(ctx, req.ShopID)
	if err != nil {
		return apperror.NewInternal(fmt.Sprintf("Failed to check shop existence: %s", err.Error()))
	}

	if !isShopExists {
		return apperror.NewNotFound("Shop", req.ShopID)
	}

	// Validate shipping address
	if req.ShippingAddressID == "" {
		return apperror.NewBadRequest("Address cannot be empty", errors.New("shipping_address_id is required"))
	}

	address, err := u.userAdapter.GetAddressById(ctx, req.ShippingAddressID)
	if err != nil {
		return apperror.NewInternal(fmt.Sprintf("Failed to get address: %s", err.Error()))
	}

	if address == nil || address.DeletedAt != nil {
		return apperror.NewNotFound("Shipping address", req.ShippingAddressID)
	}

	if len(req.Items) == 0 {
		return apperror.NewBadRequest("Order must contain at least one item", nil)
	}

	for _, item := range req.Items {
		if item.ProductID == "" {
			return apperror.NewBadRequest("Product ID cannot be empty", errors.New("product_id is required"))
		}
		if item.Quantity <= 0 {
			return apperror.NewBadRequest(fmt.Sprintf("Quantity for product %s must be positive", item.ProductID), nil)
		}
	}

	return nil
}
