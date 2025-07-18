// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package sqlc

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type PromotionStatus string

const (
	PromotionStatusDRAFT    PromotionStatus = "DRAFT"
	PromotionStatusACTIVE   PromotionStatus = "ACTIVE"
	PromotionStatusINACTIVE PromotionStatus = "INACTIVE"
	PromotionStatusDELETED  PromotionStatus = "DELETED"
)

func (e *PromotionStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PromotionStatus(s)
	case string:
		*e = PromotionStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for PromotionStatus: %T", src)
	}
	return nil
}

type NullPromotionStatus struct {
	PromotionStatus PromotionStatus `json:"promotion_status"`
	Valid           bool            `json:"valid"` // Valid is true if PromotionStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPromotionStatus) Scan(value interface{}) error {
	if value == nil {
		ns.PromotionStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PromotionStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPromotionStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PromotionStatus), nil
}

type PromotionType string

const (
	PromotionTypePERCENTAGE PromotionType = "PERCENTAGE"
	PromotionTypeVALUE      PromotionType = "VALUE"
)

func (e *PromotionType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PromotionType(s)
	case string:
		*e = PromotionType(s)
	default:
		return fmt.Errorf("unsupported scan type for PromotionType: %T", src)
	}
	return nil
}

type NullPromotionType struct {
	PromotionType PromotionType `json:"promotion_type"`
	Valid         bool          `json:"valid"` // Valid is true if PromotionType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPromotionType) Scan(value interface{}) error {
	if value == nil {
		ns.PromotionType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PromotionType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPromotionType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PromotionType), nil
}

type Address struct {
	ID        pgtype.UUID        `json:"id"`
	ShopID    pgtype.UUID        `json:"shop_id"`
	Street    string             `json:"street"`
	Ward      pgtype.Text        `json:"ward"`
	District  pgtype.Text        `json:"district"`
	City      pgtype.Text        `json:"city"`
	Country   pgtype.Text        `json:"country"`
	Lat       pgtype.Float8      `json:"lat"`
	Long      pgtype.Float8      `json:"long"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type Shop struct {
	ID              pgtype.UUID        `json:"id"`
	OwnerID         pgtype.UUID        `json:"owner_id"`
	ShopName        string             `json:"shop_name"`
	AvatarUrl       string             `json:"avatar_url"`
	BannerUrl       string             `json:"banner_url"`
	ShopDescription pgtype.Text        `json:"shop_description"`
	AddressID       pgtype.UUID        `json:"address_id"`
	Phone           string             `json:"phone"`
	Email           string             `json:"email"`
	Rating          pgtype.Numeric     `json:"rating"`
	ActiveAt        pgtype.Timestamptz `json:"active_at"`
	BannedAt        pgtype.Timestamptz `json:"banned_at"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
	UpdatedAt       pgtype.Timestamptz `json:"updated_at"`
}

type ShopPromotion struct {
	ID                pgtype.UUID         `json:"id"`
	ShopID            pgtype.UUID         `json:"shop_id"`
	PromotionName     string              `json:"promotion_name"`
	PromotionType     PromotionType       `json:"promotion_type"`
	DiscountValue     pgtype.Numeric      `json:"discount_value"`
	MaxDiscountAmount pgtype.Numeric      `json:"max_discount_amount"`
	MinPurchaseAmount pgtype.Numeric      `json:"min_purchase_amount"`
	UsageLimitPerUser pgtype.Int4         `json:"usage_limit_per_user"`
	StartTime         pgtype.Timestamptz  `json:"start_time"`
	EndTime           pgtype.Timestamptz  `json:"end_time"`
	PromotionStatus   NullPromotionStatus `json:"promotion_status"`
	CreatedAt         pgtype.Timestamptz  `json:"created_at"`
	UpdatedAt         pgtype.Timestamptz  `json:"updated_at"`
}

type ShopPromotionUsage struct {
	ID          pgtype.UUID        `json:"id"`
	PromotionID pgtype.UUID        `json:"promotion_id"`
	UserID      pgtype.UUID        `json:"user_id"`
	UsedAt      pgtype.Timestamptz `json:"used_at"`
	OrderID     pgtype.UUID        `json:"order_id"`
	UsageCount  pgtype.Int4        `json:"usage_count"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}
