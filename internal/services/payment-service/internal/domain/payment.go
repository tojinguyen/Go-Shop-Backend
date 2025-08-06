package domain

import (
	"time"

	. "github.com/toji-dev/go-shop/internal/services/payment-service/internal/constant"
)

type Payment struct {
	ID                    string        `json:"id"`
	OrderID               string        `json:"order_id"`
	UserID                string        `json:"user_id"`
	Amount                float64       `json:"amount"`
	Currency              string        `json:"currency"`
	Method                PaymentMethod `json:"payment_method"`
	Provider              string        `json:"payment_provider"`
	ProviderTransactionID *string       `json:"provider_transaction_id"`
	Status                PaymentStatus `json:"payment_status"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
}

type PaymentRefund struct {
	ID                string       `json:"id"`
	PaymentID         string       `json:"payment_id"`
	OrderID           string       `json:"order_id"`
	Amount            float64      `json:"amount"`
	Reason            string       `json:"reason"`
	ProviderPaymentID *string      `json:"provider_refund_id"`
	RefundStatus      RefundStatus `json:"refund_status"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}
