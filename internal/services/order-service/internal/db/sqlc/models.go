// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package sqlc

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type OrderStatus string

const (
	OrderStatusPENDING        OrderStatus = "PENDING"
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAYMENTFAILED  OrderStatus = "PAYMENT_FAILED"
	OrderStatusPROCESSING     OrderStatus = "PROCESSING"
	OrderStatusSHIPPED        OrderStatus = "SHIPPED"
	OrderStatusDELIVERING     OrderStatus = "DELIVERING"
	OrderStatusDELIVERED      OrderStatus = "DELIVERED"
	OrderStatusCANCELED       OrderStatus = "CANCELED"
	OrderStatusFAILED         OrderStatus = "FAILED"
)

func (e *OrderStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = OrderStatus(s)
	case string:
		*e = OrderStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for OrderStatus: %T", src)
	}
	return nil
}

type NullOrderStatus struct {
	OrderStatus OrderStatus `json:"order_status"`
	Valid       bool        `json:"valid"` // Valid is true if OrderStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullOrderStatus) Scan(value interface{}) error {
	if value == nil {
		ns.OrderStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.OrderStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullOrderStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.OrderStatus), nil
}

type Order struct {
	ID                pgtype.UUID        `json:"id"`
	OwnerID           pgtype.UUID        `json:"owner_id"`
	ShopID            pgtype.UUID        `json:"shop_id"`
	ShippingAddressID pgtype.UUID        `json:"shipping_address_id"`
	PromotionID       pgtype.UUID        `json:"promotion_id"`
	ShippingFee       pgtype.Numeric     `json:"shipping_fee"`
	DiscountAmount    pgtype.Numeric     `json:"discount_amount"`
	TotalAmount       pgtype.Numeric     `json:"total_amount"`
	FinalAmount       pgtype.Numeric     `json:"final_amount"`
	OrderStatus       OrderStatus        `json:"order_status"`
	CreatedAt         pgtype.Timestamptz `json:"created_at"`
	UpdatedAt         pgtype.Timestamptz `json:"updated_at"`
}

type OrderItem struct {
	ID        pgtype.UUID        `json:"id"`
	OrderID   pgtype.UUID        `json:"order_id"`
	ProductID pgtype.UUID        `json:"product_id"`
	ShopID    pgtype.UUID        `json:"shop_id"`
	Quantity  int32              `json:"quantity"`
	Price     pgtype.Numeric     `json:"price"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}
