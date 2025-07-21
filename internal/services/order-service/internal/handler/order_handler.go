package handler

import "github.com/toji-dev/go-shop/internal/services/order-service/internal/usecase"

type OrderHandler interface {
}

type orderHandler struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(orderUsecase usecase.OrderUsecase) OrderHandler {
	return &orderHandler{orderUsecase: orderUsecase}
}
