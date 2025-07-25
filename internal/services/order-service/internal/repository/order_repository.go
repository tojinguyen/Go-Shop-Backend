package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/pkg/converter"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/domain"
)

type OrderRepository interface {
	CreateOrder(ctx *gin.Context, order *domain.Order) (*domain.Order, error)
	UpdateOrderStatus(ctx *gin.Context, orderID string, status sqlc.OrderStatus) (*domain.Order, error)
	GetStaleOrders(ctx context.Context, olderThan time.Time, limit int) ([]*domain.Order, error)
}

type orderRepository struct {
	db      *postgresql_infra.PostgreSQLService
	queries *sqlc.Queries
}

func NewOrderRepository(db *postgresql_infra.PostgreSQLService) OrderRepository {
	if db == nil {
		return nil
	}

	queries := sqlc.New(db.GetPool())

	return &orderRepository{
		db:      db,
		queries: queries,
	}
}

func (r *orderRepository) CreateOrder(ctx *gin.Context, order *domain.Order) (*domain.Order, error) {
	tx, err := r.db.BeginTransaction(ctx)
	if err != nil {
		return nil, apperror.NewInternal(fmt.Sprintf("failed to begin transaction: %v", err))
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Step 1: Create the main order record
	orderParams := sqlc.CreateOrderParams{
		ID:                converter.StringToPgUUID(order.ID),
		UserID:            converter.StringToPgUUID(order.OwnerID),
		ShopID:            converter.StringToPgUUID(order.ShopID),
		ShippingAddressID: converter.StringToPgUUID(order.ShippingAddressID),
		OrderStatus:       sqlc.OrderStatus(order.Status),
	}
	if order.PromotionCode != nil {
		orderParams.PromotionID = converter.StringToPgUUID(*order.PromotionCode)
	}

	createdOrder, err := qtx.CreateOrder(ctx, orderParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create order record: %w", err)
	}

	// Step 2: Create order items
	createdItems := make([]domain.OrderItem, len(order.Items))
	for i, item := range order.Items {
		itemParams := sqlc.CreateOrderItemParams{
			OrderID:   createdOrder.ID,
			ProductID: converter.StringToPgUUID(item.ProductID),
			ShopID:    converter.StringToPgUUID(order.ShopID), // All items belong to the same shop
			Quantity:  int32(item.Quantity),
			Price:     converter.Float64ToPgNumeric(item.Price),
		}
		createdItem, err := qtx.CreateOrderItem(ctx, itemParams)
		if err != nil {
			return nil, fmt.Errorf("failed to create order item for product %s: %w", item.ProductID, err)
		}
		createdItems[i] = toDomainOrderItem(&createdItem)
	}

	// Step 3: Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Build the final domain object
	finalOrder := toDomainOrder(&createdOrder)
	finalOrder.Items = createdItems
	return finalOrder, nil
}

func (r *orderRepository) UpdateOrderStatus(ctx *gin.Context, orderID string, status sqlc.OrderStatus) (*domain.Order, error) {
	updatedOrder, err := r.queries.UpdateOrderStatus(ctx, sqlc.UpdateOrderStatusParams{
		ID:          converter.StringToPgUUID(orderID),
		OrderStatus: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}
	return toDomainOrder(&updatedOrder), nil
}

func (r *orderRepository) GetStaleOrders(ctx context.Context, olderThan time.Time, limit int) ([]*domain.Order, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	if olderThan.IsZero() {
		return nil, fmt.Errorf("olderThan time cannot be zero")
	}

	param := sqlc.GetStaleOrdersParams{
		UpdatedAt: converter.TimePtrToPgTime(&olderThan),
		Limit:     int32(limit),
	}

	orders, err := r.queries.GetStaleOrders(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("failed to get stale orders: %w", err)
	}

	domainOrders := make([]*domain.Order, len(orders))
	for i, order := range orders {
		domainOrders[i] = toDomainOrder(&order)
	}
	return domainOrders, nil
}

func toDomainOrder(dbOrder *sqlc.Order) *domain.Order {
	if dbOrder == nil {
		return nil
	}

	promotionCode := converter.PgUUIDToString(dbOrder.PromotionID)

	return &domain.Order{
		ID:                converter.PgUUIDToString(dbOrder.ID),
		OwnerID:           converter.PgUUIDToString(dbOrder.UserID),
		ShopID:            converter.PgUUIDToString(dbOrder.ShopID),
		ShippingAddressID: converter.PgUUIDToString(dbOrder.ShippingAddressID),
		Status:            domain.OrderStatus(dbOrder.OrderStatus),
		PromotionCode:     &promotionCode,
	}
}

func toDomainOrderItem(dbItem *sqlc.OrderItem) domain.OrderItem {
	if dbItem == nil {
		return domain.OrderItem{}
	}

	return domain.OrderItem{
		ProductID: converter.PgUUIDToString(dbItem.ProductID),
		Quantity:  int(dbItem.Quantity),
		Price:     converter.PgNumericToFloat64(dbItem.Price),
	}
}
