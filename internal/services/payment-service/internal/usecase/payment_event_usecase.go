package usecase

import (
	"context"
	"log"

	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/domain"
	grpc_adapter "github.com/toji-dev/go-shop/internal/services/payment-service/internal/grpc/adapter"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/repository"
)

const (
	BATCH_SIZE      = 100
	MAX_RETRY_COUNT = 5 // Số lần thử lại tối đa trước khi đánh dấu là FAILED
)

type PaymentEventUseCase interface {
	HandlePaymentEvent()
}

type paymentEventUseCase struct {
	eventRepo    repository.PaymentEventRepository
	orderAdapter grpc_adapter.OrderServiceAdapter
}

func NewPaymentEventUseCase(eventRepo repository.PaymentEventRepository, orderAdapter grpc_adapter.OrderServiceAdapter) PaymentEventUseCase {
	return &paymentEventUseCase{
		eventRepo:    eventRepo,
		orderAdapter: orderAdapter,
	}
}

func (uc *paymentEventUseCase) HandlePaymentEvent() {
	ctx := context.Background()
	log.Println("[PaymentEventWorker] Starting to handle pending payment events...")

	// 1. Lấy một batch các event đang ở trạng thái PENDING
	events, err := uc.eventRepo.GetBatchPaymentEventByEventTypeAndStatus(
		ctx,
		domain.PaymentEventTypePaymentSuccess,
		domain.PaymentEventStatusPending,
		BATCH_SIZE,
	)
	if err != nil {
		log.Printf("[PaymentEventWorker] Error fetching pending events: %v", err)
		return
	}

	if len(events) == 0 {
		log.Println("[PaymentEventWorker] No pending payment events to process.")
		return
	}

	log.Printf("[PaymentEventWorker] Found %d pending events to process.", len(events))
}
