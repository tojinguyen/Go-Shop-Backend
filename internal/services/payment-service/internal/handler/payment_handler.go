package handler

import "github.com/toji-dev/go-shop/internal/services/payment-service/internal/usecase"

type PaymentHandler interface {
}

type paymentHandler struct {
	orderUsecase usecase.PaymentUseCase
}

func NewPaymentHandler(orderUsecase usecase.PaymentUseCase) PaymentHandler {
	return &paymentHandler{orderUsecase: orderUsecase}
}
