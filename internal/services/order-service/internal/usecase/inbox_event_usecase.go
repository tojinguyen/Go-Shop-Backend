package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/toji-dev/go-shop/internal/services/order-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/repository"
)

const (
	INBOX_BATCH_SIZE       = 100
	INBOX_MAX_RETRY        = 5
	SOURCE_PAYMENT_SERVICE = "payment-service"
)

type InboxEventUseCase interface {
	// Kafka message handler - idempotent processing
	HandleIncomingEvent(ctx context.Context, key, value []byte) error

	// Background worker methods
	ProcessPendingInboxEvents(ctx context.Context) error
	RetryFailedInboxEvents(ctx context.Context) error
	CleanupOldEvents(ctx context.Context) error

	// Stats and monitoring
	GetInboxStats(ctx context.Context) (*repository.InboxEventStats, error)
}

type inboxEventUseCase struct {
	inboxRepo repository.InboxEventRepository
	orderRepo repository.OrderRepository
}

func NewInboxEventUseCase(
	inboxRepo repository.InboxEventRepository,
	orderRepo repository.OrderRepository,
) InboxEventUseCase {
	return &inboxEventUseCase{
		inboxRepo: inboxRepo,
		orderRepo: orderRepo,
	}
}

// HandleIncomingEvent - Main Kafka message handler with idempotency
func (uc *inboxEventUseCase) HandleIncomingEvent(ctx context.Context, key, value []byte) error {
	log.Printf("[InboxHandler] Received message with key: %s", string(key))

	// Parse the incoming event
	var kafkaPayload map[string]interface{}
	if err := json.Unmarshal(value, &kafkaPayload); err != nil {
		log.Printf("[InboxHandler] Failed to unmarshal message: %v", err)
		return fmt.Errorf("failed to unmarshal kafka message: %w", err)
	}

	// Extract event metadata
	eventID := generateEventID(key, value) // Generate deterministic ID for idempotency
	eventType := determineEventType(kafkaPayload)

	if eventType == "" {
		log.Printf("[InboxHandler] Unknown event type from payload: %v", kafkaPayload)
		return fmt.Errorf("unknown event type")
	}

	// Check if event already exists (idempotency check)
	existingEvent, err := uc.inboxRepo.GetInboxEventByEventID(ctx, eventID)
	if err == nil && existingEvent != nil {
		log.Printf("[InboxHandler] Event %s already exists with status %s", eventID, existingEvent.EventStatus)

		// If already processed, just return success
		if existingEvent.EventStatus == domain.InboxEventStatusProcessed {
			return nil
		}

		// If pending or failed, let it be processed by worker
		return nil
	}

	// Create new inbox event
	inboxEvent := &domain.InboxEvent{
		EventID:       eventID,
		EventType:     eventType,
		SourceService: SOURCE_PAYMENT_SERVICE,
		Payload:       string(value),
		EventStatus:   domain.InboxEventStatusPending,
		RetryCount:    0,
		MaxRetry:      INBOX_MAX_RETRY,
	}

	_, err = uc.inboxRepo.CreateInboxEvent(ctx, inboxEvent)
	if err != nil {
		log.Printf("[InboxHandler] Failed to create inbox event: %v", err)
		return fmt.Errorf("failed to create inbox event: %w", err)
	}

	log.Printf("[InboxHandler] Successfully created inbox event %s of type %s", eventID, eventType)
	return nil
}

// ProcessPendingInboxEvents - Background worker to process pending events
func (uc *inboxEventUseCase) ProcessPendingInboxEvents(ctx context.Context) error {
	events, err := uc.inboxRepo.GetPendingInboxEvents(ctx, INBOX_BATCH_SIZE)
	if err != nil {
		log.Printf("[InboxWorker] Failed to get pending events: %v", err)
		return err
	}

	if len(events) == 0 {
		log.Println("[InboxWorker] No pending events to process")
		return nil
	}

	log.Printf("[InboxWorker] Processing %d pending events", len(events))

	for _, event := range events {
		err := uc.processEvent(ctx, event)
		if err != nil {
			log.Printf("[InboxWorker] Failed to process event %s: %v", event.EventID, err)
			uc.markEventAsFailed(ctx, event, err)
		} else {
			uc.markEventAsProcessed(ctx, event)
		}
	}

	log.Printf("[InboxWorker] Finished processing batch of %d events", len(events))
	return nil
}

// RetryFailedInboxEvents - Retry failed events that haven't exceeded max retry
func (uc *inboxEventUseCase) RetryFailedInboxEvents(ctx context.Context) error {
	events, err := uc.inboxRepo.GetFailedInboxEvents(ctx, INBOX_BATCH_SIZE)
	if err != nil {
		log.Printf("[InboxRetryWorker] Failed to get failed events: %v", err)
		return err
	}

	if len(events) == 0 {
		log.Println("[InboxRetryWorker] No failed events to retry")
		return nil
	}

	log.Printf("[InboxRetryWorker] Retrying %d failed events", len(events))

	for _, event := range events {
		err := uc.processEvent(ctx, event)
		if err != nil {
			log.Printf("[InboxRetryWorker] Failed to retry event %s (attempt %d): %v",
				event.EventID, event.RetryCount+1, err)
			uc.markEventAsFailed(ctx, event, err)
		} else {
			uc.markEventAsProcessed(ctx, event)
		}
	}

	return nil
}

// processEvent - Core business logic for processing different event types
func (uc *inboxEventUseCase) processEvent(ctx context.Context, event *domain.InboxEvent) error {
	log.Printf("[InboxProcessor] Processing event %s of type %s", event.EventID, event.EventType)

	switch event.EventType {
	case string(domain.InboxEventTypeRefundSucceeded):
		return uc.handleRefundSucceededEvent(ctx, event)
	case string(domain.InboxEventTypePaymentSuccess):
		return uc.handlePaymentSuccessEvent(ctx, event)
	default:
		return fmt.Errorf("unsupported event type: %s", event.EventType)
	}
}

// handleRefundSucceededEvent - Process refund succeeded events
func (uc *inboxEventUseCase) handleRefundSucceededEvent(ctx context.Context, event *domain.InboxEvent) error {
	var payload domain.EventPayload
	if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal refund succeeded payload: %w", err)
	}

	if payload.OrderID == "" {
		return fmt.Errorf("order_id is required in refund succeeded event")
	}

	log.Printf("[InboxProcessor] Processing refund succeeded for OrderID: %s", payload.OrderID)

	// Update order status to REFUNDED
	_, err := uc.orderRepo.UpdateOrderStatus(ctx, payload.OrderID, sqlc.OrderStatus(domain.OrderStatusREFUNDED))
	if err != nil {
		return fmt.Errorf("failed to update order status to REFUNDED: %w", err)
	}

	log.Printf("[InboxProcessor] Successfully updated OrderID %s status to REFUNDED", payload.OrderID)
	return nil
}

// handlePaymentSuccessEvent - Process payment success events
func (uc *inboxEventUseCase) handlePaymentSuccessEvent(ctx context.Context, event *domain.InboxEvent) error {
	var payload domain.EventPayload
	if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payment success payload: %w", err)
	}

	if payload.OrderID == "" {
		return fmt.Errorf("order_id is required in payment success event")
	}

	log.Printf("[InboxProcessor] Processing payment success for OrderID: %s", payload.OrderID)

	// Update order status to PROCESSING
	_, err := uc.orderRepo.UpdateOrderStatus(ctx, payload.OrderID, sqlc.OrderStatusPROCESSING)
	if err != nil {
		return fmt.Errorf("failed to update order status to PROCESSING: %w", err)
	}

	log.Printf("[InboxProcessor] Successfully updated OrderID %s status to PROCESSING", payload.OrderID)
	return nil
}

// markEventAsProcessed - Mark event as successfully processed
func (uc *inboxEventUseCase) markEventAsProcessed(ctx context.Context, event *domain.InboxEvent) {
	event.EventStatus = domain.InboxEventStatusProcessed
	_, err := uc.inboxRepo.UpdateInboxEventStatus(ctx, event)
	if err != nil {
		log.Printf("[InboxProcessor] CRITICAL: Failed to mark event %s as processed: %v", event.EventID, err)
	} else {
		log.Printf("[InboxProcessor] Successfully marked event %s as processed", event.EventID)
	}
}

// markEventAsFailed - Mark event as failed and increment retry count
func (uc *inboxEventUseCase) markEventAsFailed(ctx context.Context, event *domain.InboxEvent, processingErr error) {
	event.RetryCount++
	if event.RetryCount >= event.MaxRetry {
		event.EventStatus = domain.InboxEventStatusFailed
		log.Printf("[InboxProcessor] CRITICAL: Event %s failed permanently after %d retries. Error: %v",
			event.EventID, event.MaxRetry, processingErr)
	}

	_, err := uc.inboxRepo.UpdateInboxEventStatus(ctx, event)
	if err != nil {
		log.Printf("[InboxProcessor] CRITICAL: Failed to update failed event %s status: %v", event.EventID, err)
	}
}

// CleanupOldEvents - Clean up old processed events
func (uc *inboxEventUseCase) CleanupOldEvents(ctx context.Context) error {
	log.Println("[InboxCleanup] Starting cleanup of old processed events")
	err := uc.inboxRepo.CleanupOldInboxEvents(ctx)
	if err != nil {
		log.Printf("[InboxCleanup] Failed to cleanup old events: %v", err)
		return err
	}
	log.Println("[InboxCleanup] Successfully cleaned up old processed events")
	return nil
}

// GetInboxStats - Get inbox statistics for monitoring
func (uc *inboxEventUseCase) GetInboxStats(ctx context.Context) (*repository.InboxEventStats, error) {
	return uc.inboxRepo.GetInboxEventStats(ctx)
}

// Helper functions

// generateEventID - Generate deterministic event ID for idempotency
func generateEventID(key, value []byte) string {
	// Try to extract event_id from payload first
	var kafkaPayload map[string]interface{}
	if err := json.Unmarshal(value, &kafkaPayload); err == nil {
		if eventID, ok := kafkaPayload["event_id"].(string); ok && eventID != "" {
			return eventID
		}
	}

	// Fallback: generate deterministic hash from key + value
	hasher := sha256.New()
	hasher.Write(key)
	hasher.Write(value)
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

// determineEventType - Determine event type from Kafka payload
func determineEventType(payload map[string]interface{}) string {
	// Check for explicit event_type field
	if eventType, ok := payload["event_type"].(string); ok {
		return eventType
	}

	// Fallback: try to infer from other fields
	if _, hasOrderID := payload["order_id"]; hasOrderID {
		if _, hasPaymentID := payload["payment_id"]; hasPaymentID {
			return string(domain.InboxEventTypeRefundSucceeded) // Default assumption
		}
	}

	return "" // Unknown event type
}
