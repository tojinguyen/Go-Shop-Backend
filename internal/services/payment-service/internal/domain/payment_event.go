package domain

import "time"

type PaymentEventStatus string

const (
	PaymentEventStatusPending PaymentEventStatus = "PENDING"
	PaymentEventStatusSent    PaymentEventStatus = "SENT"
	PaymentEventStatusFailed  PaymentEventStatus = "FAILED"
)

type PaymentEvent struct {
	ID          string             `json:"id"`
	PaymentID   string             `json:"payment_id"`
	EventType   string             `json:"event_type"`
	Payload     string             `json:"payload"`
	EventStatus PaymentEventStatus `json:"event_status"`
	RetryCount  int                `json:"retry_count"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}
