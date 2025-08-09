package domain

import "time"

type PaymentEventStatus string

const (
	PaymentEventStatusPending PaymentEventStatus = "PENDING"
	PaymentEventStatusSent    PaymentEventStatus = "SENT"
	PaymentEventStatusFailed  PaymentEventStatus = "FAILED"
)

type PaymentEventType string

const (
	PaymentEventTypePaymentSuccess  PaymentEventType = "PAYMENT_SUCCESS"
	PaymentEventTypePaymentFailed   PaymentEventType = "PAYMENT_FAILED"
	PaymentEventTypeRefundRequested PaymentEventType = "REFUND_REQUESTED"
	PaymentEventTypeRefundSuccessed PaymentEventType = "REFUND_SUCCEEDED"
)

type PaymentEvent struct {
	ID          string             `json:"id"`
	PaymentID   string             `json:"payment_id"`
	OrderID     string             `json:"order_id"`
	EventType   string             `json:"event_type"`
	Payload     string             `json:"payload"`
	EventStatus PaymentEventStatus `json:"event_status"`
	RetryCount  int                `json:"retry_count"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}
