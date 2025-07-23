package domain

type Order struct {
	ID                string      `json:"id"`
	OwnerID           string      `json:"customer_id"`
	ShopID            string      `json:"shop_id"`
	ShippingAddressID string      `json:"shipping_address_id"`
	PromotionCode     *string     `json:"promotion_code,omitempty"`
	DiscountAmount    float64     `json:"discount_amount"`
	TotalAmount       float64     `json:"total_amount"`
	Status            string      `json:"status"`
	Items             []OrderItem `json:"items"`
	CreatedAt         string      `json:"created_at"`
	UpdatedAt         string      `json:"updated_at"`
}
