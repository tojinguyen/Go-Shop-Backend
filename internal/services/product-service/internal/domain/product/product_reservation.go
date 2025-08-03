package product

import "time"

type ProductReservationStatus string

const (
	ProductReservationStatusReserved   ProductReservationStatus = "RESERVED"
	ProductReservationStatusUnreserved ProductReservationStatus = "UNRESERVED"
)

type ProductReservation struct {
	ID        string                   `json:"id"`
	OrderID   string                   `json:"order_id"`
	ProductID string                   `json:"product_id"`
	ShopID    string                   `json:"shop_id"`
	Status    ProductReservationStatus `json:"status"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}
