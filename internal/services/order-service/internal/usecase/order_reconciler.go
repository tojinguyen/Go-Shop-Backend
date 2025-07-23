package usecase

import (
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/grpc/adapter"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/repository"
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
}
