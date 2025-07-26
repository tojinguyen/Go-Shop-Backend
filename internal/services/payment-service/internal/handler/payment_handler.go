package handler

import (
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/usecase"
)

type PaymentHandler interface {
}

type paymentHandler struct {
	orderUsecase usecase.PaymentUseCase
}

func NewOrderHandler(orderUsecase usecase.PaymentUseCase) PaymentHandler {
	return &paymentHandler{orderUsecase: orderUsecase}
}
