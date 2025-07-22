package dto

type CreateOrderRequest struct {
	CartID string `json:"cart_id" validate:"required"`
}
