package dto

import "time"

type CreateAddressRequest struct {
	IsDefault bool     `json:"is_default"`
	Street    string   `json:"street" validate:"required,min=1,max=255"`
	Ward      *string  `json:"ward,omitempty" validate:"omitempty,max=100"`
	District  *string  `json:"district,omitempty" validate:"omitempty,max=100"`
	City      *string  `json:"city,omitempty" validate:"omitempty,max=100"`
	Country   string   `json:"country" validate:"max=100"`
	Lat       *float64 `json:"lat,omitempty"`
	Long      *float64 `json:"long,omitempty"`
}

type UpdateAddressRequest struct {
	IsDefault bool     `json:"is_default"`
	Street    string   `json:"street" validate:"required,min=1,max=255"`
	Ward      *string  `json:"ward,omitempty" validate:"omitempty,max=100"`
	District  *string  `json:"district,omitempty" validate:"omitempty,max=100"`
	City      *string  `json:"city,omitempty" validate:"omitempty,max=100"`
	Country   string   `json:"country" validate:"max=100"`
	Lat       *float64 `json:"lat,omitempty"`
	Long      *float64 `json:"long,omitempty"`
}

type GetAddressByUserIDRequest struct {
	UserID string `uri:"user_id" validate:"required,uuid"`
}

type GetDefaultAddressByUserIDRequest struct {
	UserID string `uri:"user_id" validate:"required,uuid"`
}

type SetDefaultAddressRequest struct {
	UserID    string `json:"user_id" validate:"required,uuid"`
	AddressID string `json:"address_id" validate:"required,uuid"`
}

type AddressResponse struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	IsDefault bool       `json:"is_default"`
	Street    string     `json:"street"`
	Ward      *string    `json:"ward,omitempty"`
	District  *string    `json:"district,omitempty"`
	City      *string    `json:"city,omitempty"`
	Country   string     `json:"country"`
	Lat       *float64   `json:"lat,omitempty"`
	Long      *float64   `json:"long,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type AddressListResponse struct {
	Addresses []AddressResponse `json:"addresses"`
	Total     int               `json:"total"`
}
