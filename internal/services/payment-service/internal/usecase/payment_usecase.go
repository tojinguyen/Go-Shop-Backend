package usecase

import (
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/repository"
)

type PaymentUseCase interface {
}

type paymentUseCase struct {
	paymentRepo repository.PaymentRepository
}

func NewPaymentUsecase(paymentRepo repository.PaymentRepository) PaymentUseCase {
	return &paymentUseCase{paymentRepo: paymentRepo}
}
