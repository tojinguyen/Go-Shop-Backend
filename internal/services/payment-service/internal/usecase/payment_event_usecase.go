package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/domain"
	grpc_adapter "github.com/toji-dev/go-shop/internal/services/payment-service/internal/grpc/adapter"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/repository"
	order_v1 "github.com/toji-dev/go-shop/proto/gen/go/order/v1"
)

const (
	BATCH_SIZE      = 100
	MAX_RETRY_COUNT = 5 // Số lần thử lại tối đa trước khi đánh dấu là FAILED
)

type PaymentEventUseCase interface {
	HandleSuccessPaymentPending()
	HandleRefundPaymentPending()
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

func (uc *paymentEventUseCase) HandleSuccessPaymentPending() {
	ctx := context.Background()
	log.Println("[PaymentEventWorker] Starting to handle pending payment events...")

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

	for _, event := range events {
		err := uc.processEvent(ctx, event)

		if err != nil {
			log.Printf("[PaymentEventWorker] Error processing event ID %s: %v. Updating retry count.", event.ID, err)
			event.RetryCount++
			if event.RetryCount >= MAX_RETRY_COUNT {
				event.EventStatus = domain.PaymentEventStatusFailed
			}
			_, updateErr := uc.eventRepo.UpdatePaymentEvent(ctx, event)
			if updateErr != nil {
				log.Printf("[PaymentEventWorker] CRITICAL: Failed to update event status after processing error for event ID %s: %v", event.ID, updateErr)
			}
		} else {
			event.EventStatus = domain.PaymentEventStatusSent
			_, updateErr := uc.eventRepo.UpdatePaymentEvent(ctx, event)
			if updateErr != nil {
				log.Printf("[PaymentEventWorker] CRITICAL: Failed to update event status after successful processing for event ID %s: %v", event.ID, updateErr)
			}
			log.Printf("[PaymentEventWorker] Successfully processed event ID %s for Order ID %s.", event.ID, event.OrderID)
		}
	}

	log.Printf("[PaymentEventWorker] Finished processing batch of %d events.", len(events))
}

func (uc *paymentEventUseCase) HandleRefundPaymentPending() {
}

func (uc *paymentEventUseCase) processEvent(ctx context.Context, event *domain.PaymentEvent) error {
	var payload struct {
		PaymentStatus string `json:"payment_status"`
	}

	if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal event payload: %w", err)
	}

	var newOrderStatus order_v1.OrderStatus
	switch payload.PaymentStatus {
	case string(sqlc.PaymentStatusSUCCESS):
		newOrderStatus = order_v1.OrderStatus_ORDER_STATUS_PROCESSING
	case string(sqlc.PaymentStatusFAILED):
		newOrderStatus = order_v1.OrderStatus_ORDER_STATUS_PAYMENT_FAILED
	default:
		return fmt.Errorf("unhandled payment status in event: %s", payload.PaymentStatus)
	}

	req := &order_v1.UpdateOrderStatusRequest{
		OrderId:   event.OrderID,
		NewStatus: newOrderStatus,
	}

	_, err := uc.orderAdapter.UpdateOrderStatus(ctx, req)
	if err != nil {
		return fmt.Errorf("gRPC call to OrderService failed: %w", err)
	}
	return nil
}
