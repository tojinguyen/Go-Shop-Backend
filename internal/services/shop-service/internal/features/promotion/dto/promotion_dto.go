package dto

import "time"

type PromotionResponse struct {
	ID                string    `json:"id"`
	ShopID            string    `json:"shop_id"`
	PromotionName     string    `json:"promotion_name"`
	PromotionType     string    `json:"promotion_type"`
	DiscountValue     float64   `json:"discount_value"`
	MaxDiscountAmount *float64  `json:"max_discount_amount,omitempty"`
	MinPurchaseAmount float64   `json:"min_purchase_amount"`
	UsageLimitPerUser int32     `json:"usage_limit_per_user"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	PromotionStatus   string    `json:"promotion_status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
