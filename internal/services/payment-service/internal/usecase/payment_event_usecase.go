package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/toji-dev/go-shop/internal/pkg/converter"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/constant"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/domain"
	grpc_adapter "github.com/toji-dev/go-shop/internal/services/payment-service/internal/grpc/adapter"
	paymentprovider "github.com/toji-dev/go-shop/internal/services/payment-service/internal/payment_provider"
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
	paymentRepo     repository.PaymentRepository
	eventRepo       repository.PaymentEventRepository
	orderAdapter    grpc_adapter.OrderServiceAdapter
	providerFactory *paymentprovider.PaymentProviderFactory
}

func NewPaymentEventUseCase(
	paymentRepo repository.PaymentRepository,
	eventRepo repository.PaymentEventRepository,
	orderAdapter grpc_adapter.OrderServiceAdapter,
	providerFactory *paymentprovider.PaymentProviderFactory,
) PaymentEventUseCase {
	return &paymentEventUseCase{
		paymentRepo:     paymentRepo,
		eventRepo:       eventRepo,
		orderAdapter:    orderAdapter,
		providerFactory: providerFactory,
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
		err := uc.processUpdateOrderStatus(ctx, event)

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
	ctx := context.Background()
	log.Println("[PaymentEventWorker] Starting to handle pending payment events...")

	events, err := uc.eventRepo.GetBatchPaymentEventByEventTypeAndStatus(
		ctx,
		domain.PaymentEventTypeRefundRequested,
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
		err := uc.processRefundRequest(ctx, event)

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
				continue
			}

			// Add payment event REFUND_SUCCEEDED
			eventRefundSuccessed := &domain.PaymentEvent{
				PaymentID:   event.PaymentID,
				OrderID:     event.OrderID,
				EventType:   string(domain.PaymentEventTypeRefundSuccessed),
				Payload:     event.Payload,
				EventStatus: domain.PaymentEventStatusPending,
				RetryCount:  0,
			}
			_, insertErr := uc.eventRepo.CreatePaymentEvent(ctx, eventRefundSuccessed)
			if insertErr != nil {
				log.Printf("[PaymentEventWorker] CRITICAL: Failed to insert REFUND_SUCCEEDED event for event ID %s: %v", event.ID, insertErr)
			}
			log.Printf("[PaymentEventWorker] Successfully processed event ID %s for Order ID %s.", event.ID, event.OrderID)
		}
	}

	log.Printf("[PaymentEventWorker] Finished processing batch of %d events.", len(events))
}

func (uc *paymentEventUseCase) processUpdateOrderStatus(ctx context.Context, event *domain.PaymentEvent) error {
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
		log.Printf("gRPC call to OrderService failed: %v", err)
		return fmt.Errorf("gRPC call to OrderService failed: %w", err)
	}
	return nil
}

func (uc *paymentEventUseCase) processRefundRequest(ctx context.Context, event *domain.PaymentEvent) error {
	var payload struct {
		Reason string `json:"reason"`
	}

	if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal event payload: %w", err)
	}

	log.Printf("Processing refund request for OrderID %s with reason: %s", event.OrderID, payload.Reason)

	payment, err := uc.paymentRepo.GetPaymentByOrderID(ctx, event.OrderID)
	if err != nil {
		log.Printf("Failed to get payment by OrderID %s: %v", event.OrderID, err)
		return fmt.Errorf("failed to get payment by OrderID: %w", err)
	}

	paymentMethod := payment.Provider
	paymentProvider, err := uc.providerFactory.GetProvider(constant.PaymentProviderMethod(paymentMethod))

	refundPayment, err := uc.paymentRepo.GetRefundByPaymentID(ctx, payment.ID)

	// Call provider's refund method
	refundData := paymentprovider.RefundData{
		PaymentID:             event.PaymentID,
		OrderID:               event.OrderID,
		ProviderTransactionID: *payment.ProviderTransactionID,
		Amount:                int64(payment.Amount),
		Reason:                refundPayment.Reason,
	}

	refundRes, err := paymentProvider.Refund(ctx, refundData)

	if err != nil {
		log.Printf("Error refunding payment with ID %s: %v", event.PaymentID, err)
		updateRefundStatusParams := sqlc.UpdateRefundPaymentStatusParams{
			ID:           converter.StringToPgUUID(refundPayment.ID),
			RefundStatus: sqlc.RefundStatusFAILED,
		}

		_, err = uc.paymentRepo.UpdateRefundPaymentStatus(ctx, updateRefundStatusParams)
		if err != nil {
			log.Printf("Failed to update refund payment status for OrderID %s: %v", event.OrderID, err)
			return fmt.Errorf("failed to update refund payment status for OrderID %s: %w", event.OrderID, err)
		}

		return fmt.Errorf("failed to refund payment with ID %s: %w", event.PaymentID, err)
	}

	if refundRes == nil {
		log.Printf("Refund result is nil for OrderID %s", event.OrderID)
		return fmt.Errorf("refund result is nil for OrderID %s", event.OrderID)
	}

	updateRefundStatusParams := sqlc.UpdateRefundPaymentStatusParams{
		ID:           converter.StringToPgUUID(refundPayment.ID),
		RefundStatus: sqlc.RefundStatusCOMPLETED,
	}

	_, err = uc.paymentRepo.UpdateRefundPaymentStatus(ctx, updateRefundStatusParams)
	if err != nil {
		log.Printf("Failed to update refund payment status for OrderID %s: %v", event.OrderID, err)
		return fmt.Errorf("failed to update refund payment status for OrderID %s: %w", event.OrderID, err)
	}

	updatePaymentStatus := sqlc.UpdatePaymentStatusParams{
		ID:                    converter.StringToPgUUID(payment.ID),
		PaymentStatus:         sqlc.PaymentStatusREFUNDED,
		ProviderTransactionID: converter.StringToPgText(payment.ProviderTransactionID),
	}

	_, err = uc.paymentRepo.UpdatePaymentStatus(ctx, updatePaymentStatus)

	if err != nil {
		log.Printf("Failed to update payment status for OrderID %s: %v", event.OrderID, err)
		return fmt.Errorf("failed to update payment status for OrderID %s: %w", event.OrderID, err)
	}

	log.Printf("Refund request for OrderID %s processed successfully.", event.OrderID)

	return nil
}
