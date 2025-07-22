package domain

type Order struct {
	ID          string      `json:"id"`
	OwnerID     string      `json:"customer_id"`
	ShopID      string      `json:"shop_id"`
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"`
	Items       []OrderItem `json:"items"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}
