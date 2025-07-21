package usecase

import "github.com/toji-dev/go-shop/internal/services/order-service/internal/repository"

type OrderItemUsecase interface {
}

type orderItemUsecase struct {
	orderRepo repository.OrderRepository
}

func NewOrderItemUsecase(orderRepo repository.OrderRepository) OrderItemUsecase {
	return &orderItemUsecase{orderRepo: orderRepo}
}
