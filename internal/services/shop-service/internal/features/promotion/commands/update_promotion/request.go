package updatepromotion

import "time"

type UpdatePromotionRequest struct {
	PromotionName     string    `json:"promotion_name" binding:"required"`
	PromotionType     string    `json:"promotion_type" binding:"required,oneof=PERCENTAGE VALUE"`
	DiscountValue     float64   `json:"discount_value" binding:"required,gt=0"`
	MaxDiscountAmount *float64  `json:"max_discount_amount,omitempty"`
	MinPurchaseAmount float64   `json:"min_purchase_amount" binding:"gte=0"`
	UsageLimitPerUser int32     `json:"usage_limit_per_user" binding:"gte=1"`
	StartTime         time.Time `json:"start_time" binding:"required"`
	EndTime           time.Time `json:"end_time" binding:"required,gtfield=StartTime"`
	PromotionStatus   string    `json:"promotion_status" binding:"required,oneof=DRAFT ACTIVE INACTIVE"`
}
