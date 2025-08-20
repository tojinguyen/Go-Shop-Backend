package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/grpc/adapter"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/payload"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/repository"
	product_v1 "github.com/toji-dev/go-shop/proto/gen/go/product/v1"
	shop_v1 "github.com/toji-dev/go-shop/proto/gen/go/shop/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, userId string, req dto.CreateOrderRequest) (*domain.Order, error)
	HandleRefundSucceededEvent(ctx context.Context, key, value []byte) error // Deprecated: Use InboxEventUseCase instead
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

func (u *orderUsecase) CreateOrder(ctx context.Context, userId string, req dto.CreateOrderRequest) (*domain.Order, error) {
	tracer := otel.Tracer("order-service.usecase")
	ctx, span := tracer.Start(ctx, "CreateOrder.UseCase")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.id", userId),
		attribute.String("shop.id", req.ShopID),
		attribute.Int("order.item_count", len(req.Items)),
	)

	// --- STAGE 1: VALIDATION ---
	// Validate shop, address, and product info before creating anything
	_, validationSpan := tracer.Start(ctx, "ValidatePrerequisites")
	if err := u.validatePrerequisites(ctx, &req); err != nil {
		validationSpan.SetStatus(codes.Error, err.Error()) // Ghi nhận lỗi vào span
		validationSpan.End()
		return nil, err
	}
	validationSpan.End()

	log.Printf("Order request: %+v\n", req)
	log.Printf("User ID: %d", len(req.Items))

	// Validate products in order items
	productIDs := make([]string, 0, len(req.Items))
	quantityMap := make(map[string]int32)

	log.Printf("Validating products for order items")

	for i, item := range req.Items {
		log.Printf("Validating item %d: ProductID=%s, Quantity=%d", i+1, item.ProductID, item.Quantity)
		productIDs = append(productIDs, item.ProductID)
		quantityMap[item.ProductID] = int32(item.Quantity)
	}

	_, productInfoSpan := tracer.Start(ctx, "GetProductsInfo_gRPC")
	productsInfo, err := u.productServiceAdapter.GetProductsInfo(ctx, &product_v1.GetProductsInfoRequest{
		ProductIds: productIDs,
	})

	if err != nil || !productsInfo.Valid || productsInfo.Products == nil || len(productsInfo.Products) != len(req.Items) {
		productInfoSpan.SetStatus(codes.Error, "One or more products are invalid")
		productInfoSpan.End()
		return nil, apperror.NewBadRequest("One or more products are invalid or unavailable", err)
	}
	productInfoSpan.End()

	// --- STAGE 2: CALCULATION & ORDER CREATION (order_status: PENDING) ---
	_, calculationSpan := tracer.Start(ctx, "CalculateTotalsAndPromotions")
	orderID := uuid.New().String()
	totalAmount := 0.0
	orderItems := make([]domain.OrderItem, len(productsInfo.Products))

	log.Printf("Product info retrieved with %d products", len(productsInfo.Products))

	for i, p := range productsInfo.Products {
		price := float64(p.Price)
		totalAmount += price * float64(quantityMap[p.Id])
		orderItems[i] = domain.OrderItem{
			ProductID: p.Id,
			Quantity:  int(quantityMap[p.Id]),
			Price:     price,
		}
	}

	finalPrice := totalAmount

	log.Printf("Total amount calculated: %.2f", totalAmount)

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
			log.Printf("Error calculating promotion: %v", err)
			return nil, apperror.NewInternal(fmt.Sprintf("Failed to calculate promotion: %s", err.Error()))
		}

		if !promotionRes.Eligible {
			return nil, apperror.NewBadRequest("Invalid promotion code", errors.New("promotion code is not eligible"))
		}

		finalPrice -= float64(promotionRes.Discount)
		discountAmount = float64(promotionRes.Discount)
	}

	log.Printf("Final price after promotion: %.2f (Discount: %.2f)", finalPrice, discountAmount)

	calculationSpan.SetAttributes(
		attribute.Float64("order.total_amount", totalAmount),
		attribute.Float64("order.discount_amount", discountAmount),
		attribute.Float64("order.final_price", finalPrice),
	)
	calculationSpan.End()

	// Create the order with PENDING status FIRST
	pendingOrder := &domain.Order{
		ID:                orderID,
		OwnerID:           userId,
		ShopID:            req.ShopID,
		ShippingAddressID: req.ShippingAddressID,
		PromotionCode:     req.PromotionID,
		ShippingFee:       0,
		DiscountAmount:    discountAmount,
		TotalAmount:       totalAmount,
		FinalPrice:        finalPrice,
		Status:            domain.OrderStatusPENDING,
		Items:             orderItems,
	}

	createdOrder, err := u.orderRepo.CreateOrder(ctx, pendingOrder)

	if err != nil {
		span.SetStatus(codes.Error, "Failed to create initial order record")
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

	_, reserveSpan := tracer.Start(ctx, "ReserveProducts_gRPC")

	reserveResp, err := u.productServiceAdapter.ReserveProducts(ctx, reserveReq)

	if err != nil {
		reserveSpan.SetStatus(codes.Error, err.Error())
		reserveSpan.End()
		log.Printf("CRITICAL: ReserveProducts call failed for order %s. Marking as FAILED. Error: %v", orderID, err)
		u.orderRepo.UpdateOrderStatus(ctx, orderID, sqlc.OrderStatusFAILED)
		return nil, apperror.NewDependencyFailure(fmt.Sprintf("Failed to reserve products: %s", err.Error()))
	}
	reserveSpan.End()

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

	span.AddEvent("Order successfully created and stock reserved")
	return finalOrder, nil
}

func (u *orderUsecase) HandleRefundSucceededEvent(ctx context.Context, key, value []byte) error {
	var payload payload.RefundSucceededPayload
	if err := json.Unmarshal(value, &payload); err != nil {
		log.Printf("ERROR: Failed to unmarshal refund event payload: %v", err)
		return nil
	}

	if _, err := uuid.Parse(payload.OrderID); err != nil {
		log.Printf("ERROR: [POISON PILL] OrderID không hợp lệ. Bỏ qua message. OrderID: %s. Lỗi: %v", payload.OrderID, err)
		return nil
	}

	log.Printf("Received RefundSucceeded event for OrderID: %s", payload.OrderID)

	// Cập nhật trạng thái đơn hàng thành REFUNDED
	_, err := u.orderRepo.UpdateOrderStatus(ctx, payload.OrderID, sqlc.OrderStatus(domain.OrderStatusREFUNDED))
	if err != nil {
		log.Printf("ERROR: Failed to update order status to REFUNDED for OrderID %s: %v", payload.OrderID, err)
		return fmt.Errorf("failed to update order status: %w", err)
	}

	log.Printf("Successfully updated OrderID %s status to REFUNDED", payload.OrderID)
	return nil
}

// Validate data with business rules before creating an order
func (u *orderUsecase) validatePrerequisites(ctx context.Context, req *dto.CreateOrderRequest) error {
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
		log.Printf("Order creation failed: Shipping address ID is required")
		return apperror.NewBadRequest("Address cannot be empty", errors.New("shipping_address_id is required"))
	}

	address, err := u.userAdapter.GetAddressById(ctx, req.ShippingAddressID)
	if err != nil {
		log.Printf("Error fetching address: %v", err)
		return apperror.NewInternal(fmt.Sprintf("Failed to get address: %s", err.Error()))
	}

	if address == nil {
		log.Printf("Shipping address with ID %s not found or deleted with data: %+v", req.ShippingAddressID, address)
		return apperror.NewNotFound("Shipping address", req.ShippingAddressID)
	}

	if len(req.Items) == 0 {
		log.Printf("Order creation failed: No items provided")
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
