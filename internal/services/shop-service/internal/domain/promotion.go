package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PromotionType string

const (
	PromotionTypePercentage PromotionType = "PERCENTAGE"
	PromotionTypeValue      PromotionType = "VALUE"
)

type PromotionStatus string

const (
	PromotionStatusDraft    PromotionStatus = "DRAFT"
	PromotionStatusActive   PromotionStatus = "ACTIVE"
	PromotionStatusInactive PromotionStatus = "INACTIVE"
	PromotionStatusDeleted  PromotionStatus = "DELETED"
)

type Promotion struct {
	ID                uuid.UUID       `json:"id"`
	ShopID            uuid.UUID       `json:"shop_id"`
	PromotionName     string          `json:"promotion_name"`
	PromotionType     PromotionType   `json:"promotion_type"`
	DiscountValue     float64         `json:"discount_value"`
	MaxDiscountAmount *float64        `json:"max_discount_amount,omitempty"`
	MinPurchaseAmount float64         `json:"min_purchase_amount"`
	UsageLimitPerUser int32           `json:"usage_limit_per_user"`
	StartTime         time.Time       `json:"start_time"`
	EndTime           time.Time       `json:"end_time"`
	PromotionStatus   PromotionStatus `json:"promotion_status"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

func (p *Promotion) CalculateDiscount(amount float64) (float64, error) {
	if p.PromotionStatus != PromotionStatusActive {
		return 0, fmt.Errorf("promotion is not active")
	}
	if amount < p.MinPurchaseAmount {
		return 0, fmt.Errorf("amount is less than minimum purchase amount")
	}

	if p.StartTime.After(time.Now()) || p.EndTime.Before(time.Now()) {
		return 0, fmt.Errorf("promotion is not active")
	}

	discount := 0.0
	switch p.PromotionType {
	case PromotionTypePercentage:
		discount = amount * (p.DiscountValue / 100)
	case PromotionTypeValue:
		discount = p.DiscountValue
	}

	if p.MaxDiscountAmount != nil && discount > *p.MaxDiscountAmount {
		return *p.MaxDiscountAmount, nil
	}
	return discount, nil
}
