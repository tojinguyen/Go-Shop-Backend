package domain

import "time"

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "PENDING"
	PaymentStatusProcessing PaymentStatus = "PROCESSING"
	PaymentStatusCompleted  PaymentStatus = "SUCCESS"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusRefunded   PaymentStatus = "REFUNDED"
)

type PaymentMethod string

const (
	PaymentMethodCOD          PaymentMethod = "COD"
	PaymentMethodCreditCard   PaymentMethod = "CREDIT_CARD"
	PaymentMethodBankTransfer PaymentMethod = "BANK_TRANSFER"
	PaymentMethodEWallet      PaymentMethod = "E_WALLET"
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
