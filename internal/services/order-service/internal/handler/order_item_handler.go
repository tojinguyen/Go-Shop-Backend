package handler

import "github.com/toji-dev/go-shop/internal/services/order-service/internal/usecase"

type OrderItemHandler interface {
}

type orderItemHandler struct {
	orderItemUsecase usecase.OrderItemUsecase
}

func NewOrderItemHandler(orderItemUsecase usecase.OrderItemUsecase) OrderItemHandler {
	return &orderItemHandler{orderItemUsecase: orderItemUsecase}
}
