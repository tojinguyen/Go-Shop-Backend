package dto

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type CreateOrderRequest struct {
	ShopID            string                   `json:"shop_id" binding:"required,uuid"`
	ShippingAddressID string                   `json:"shipping_address_id" binding:"required,uuid"`
	PromotionID       *string                  `json:"promotion_id,omitempty" binding:"omitempty,uuid"`
	Items             []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

type OrderResponse struct {
	ID                string              `json:"id"`
	ShopID            string              `json:"shop_id"`
	ShippingAddressID string              `json:"shipping_address_id"`
	PromotionID       *string             `json:"promotion_id,omitempty"`
	ShippingFee       float64             `json:"shipping_fee"`
	DiscountAmount    float64             `json:"discount_amount"`
	TotalAmount       float64             `json:"total_amount"`
	FinalAmount       float64             `json:"final_amount"`
	Status            string              `json:"status"`
	CreatedAt         string              `json:"created_at"`
	UpdatedAt         string              `json:"updated_at"`
	Items             []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
