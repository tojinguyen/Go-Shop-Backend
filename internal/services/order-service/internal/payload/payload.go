package payload

type RefundSucceededPayload struct {
	OrderID   string `json:"order_id"`
	PaymentID string `json:"payment_id"`
}
