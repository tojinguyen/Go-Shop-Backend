package usecase

import (
	"context"
	"log"
	"time"

	time_utils "github.com/toji-dev/go-shop/internal/pkg/time"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/grpc/adapter"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/repository"
)

const (
	STALE_ORDER_THRESHOLD = 10 * time.Minute
	BATCH_SIZE            = 100
)

type OrderReconciler struct {
	orderRepo      repository.OrderRepository
	productAdapter adapter.ProductServiceAdapter
}

func NewOrderReconciler(orderRepo repository.OrderRepository, productAdapter adapter.ProductServiceAdapter) *OrderReconciler {
	return &OrderReconciler{
		orderRepo:      orderRepo,
		productAdapter: productAdapter,
	}
}

func (r *OrderReconciler) ReconcilePendingOrders() {
	ctx := context.Background()

	log.Println("[OrderReconciler] Starting reconciliation of pending orders...")

	olderThan := time_utils.GetUtcTime().Add(-STALE_ORDER_THRESHOLD)
	orders, err := r.orderRepo.GetStaleOrders(ctx, olderThan, BATCH_SIZE)

	if err != nil {
		log.Printf("[OrderReconciler] Error fetching stale orders: %v", err)
		return
	}

	if len(orders) == 0 {
		log.Println("[OrderReconciler] No stale orders found.")
		return
	}

	log.Printf("[OrderReconciler] Found %d stale orders to process.", len(orders))

	for _, order := range orders {
		log.Printf("[Worker] Processing order ID: %s", order.ID)

	}
}
