package dto

import "time"

type ShipperRegisterRequest struct {
	VehicleType     string `json:"vehicle_type,omitempty" validate:"omitempty,max=100"`
	VehicleImageURL string `json:"vehicle_image_url,omitempty" validate:"omitempty,url,max=500"`
	IdentifyCardURL string `json:"identify_card_url,omitempty" validate:"omitempty,url,max=500"`
	LicensePlate    string `json:"license_plate,omitempty" validate:"omitempty,max=20"`
}

type ShipperUpdateRequest struct {
	VehicleType     string `json:"vehicle_type,omitempty" validate:"omitempty,max=100"`
	VehicleImageURL string `json:"vehicle_image_url,omitempty" validate:"omitempty,url,max=500"`
	IdentifyCardURL string `json:"identify_card_url,omitempty" validate:"omitempty,url,max=500"`
	LicensePlate    string `json:"license_plate,omitempty" validate:"omitempty,max=20"`
}

type ShipperResponse struct {
	UserID          string    `json:"user_id"`
	VehicleType     string    `json:"vehicle_type,omitempty"`
	VehicleImageURL string    `json:"vehicle_image_url,omitempty"`
	IdentifyCardURL string    `json:"identify_card_url,omitempty"`
	LicensePlate    string    `json:"license_plate,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
