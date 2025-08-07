package domain

import "time"

// InboxEventStatus represents the status of an inbox event
type InboxEventStatus string

const (
	InboxEventStatusPending   InboxEventStatus = "PENDING"
	InboxEventStatusProcessed InboxEventStatus = "PROCESSED"
	InboxEventStatusFailed    InboxEventStatus = "FAILED"
)

// InboxEventType represents different types of events that can be received
type InboxEventType string

const (
	InboxEventTypePaymentSuccess  InboxEventType = "PAYMENT_SUCCESS"
	InboxEventTypeRefundSucceeded InboxEventType = "REFUND_SUCCEEDED"
	InboxEventTypeRefundRequested InboxEventType = "REFUND_REQUESTED"
)

// InboxEvent represents an event received from external services
type InboxEvent struct {
	ID            string           `json:"id"`
	EventID       string           `json:"event_id"`       // UUID từ hệ thống gửi
	EventType     string           `json:"event_type"`     // Loại event
	SourceService string           `json:"source_service"` // Service gửi event
	Payload       string           `json:"payload"`        // Dữ liệu event dạng JSON
	EventStatus   InboxEventStatus `json:"event_status"`   // Trạng thái xử lý
	RetryCount    int              `json:"retry_count"`    // Số lần retry
	MaxRetry      int              `json:"max_retry"`      // Số lần retry tối đa
	ReceivedAt    time.Time        `json:"received_at"`    // Thời gian nhận event
	ProcessedAt   *time.Time       `json:"processed_at"`   // Thời gian xử lý xong
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

// EventPayload represents the common structure of event payloads
type EventPayload struct {
	OrderID     string            `json:"order_id"`
	PaymentID   string            `json:"payment_id,omitempty"`
	EventSource string            `json:"event_source,omitempty"`
	Data        map[string]string `json:"data,omitempty"`
}
