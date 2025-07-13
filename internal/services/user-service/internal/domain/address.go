package domain

import "time"

type Address struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	IsDefault bool      `json:"is_default" db:"is_default"`
	Street    string    `json:"street" db:"street"`
	Ward      string    `json:"ward,omitempty" db:"ward"`
	District  string    `json:"district,omitempty" db:"district"`
	City      string    `json:"city,omitempty" db:"city"`
	Country   string    `json:"country" db:"country"`
	Lat       float64   `json:"lat,omitempty" db:"lat"`
	Long      float64   `json:"long,omitempty" db:"long"`
	DeletedAt time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
