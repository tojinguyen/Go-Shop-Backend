package domain

import (
	"time"
)

// Shipper represents a shipper profile in the system
type Shipper struct {
	UserID          string    `json:"user_id" db:"user_id"`
	VehicleType     *string   `json:"vehicle_type,omitempty" db:"vehicle_type"`
	VehicleImageURL *string   `json:"vehicle_image_url,omitempty" db:"vehicle_image_url"`
	IdentifyCardURL *string   `json:"identify_card_url,omitempty" db:"identify_card_url"`
	LicensePlate    *string   `json:"license_plate,omitempty" db:"license_plate"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
