package usecase

import "github.com/toji-dev/go-shop/internal/services/order-service/internal/repository"

type OrderUsecase interface {
}

type orderUsecase struct {
	orderRepo repository.OrderRepository
}

func NewOrderUsecase(orderRepo repository.OrderRepository) OrderUsecase {
	return &orderUsecase{orderRepo: orderRepo}
}
