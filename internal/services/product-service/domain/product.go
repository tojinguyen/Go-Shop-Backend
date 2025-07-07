package domain

import (
	"time"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
)

type ProductStatus string

const (
	ProductStatusDraft        ProductStatus = "DRAFT"
	ProductStatusActive       ProductStatus = "ACTIVE"
	ProductStatusInactive     ProductStatus = "INACTIVE"
	ProductStatusOutOfStock   ProductStatus = "OUT_OF_STOCK"
	ProductStatusDiscontinued ProductStatus = "DISCONTINUED"
	ProductStatusRejected     ProductStatus = "REJECTED"
	ProductStatusBanned       ProductStatus = "BANNED"
)

type Price struct {
	Amount   float64
	Currency string
}

type Product struct {
	ID              uuid.UUID
	ShopID          uuid.UUID
	Name            string
	ThumbnailURL    string
	Description     string
	CategoryID      *uuid.UUID
	Price           Price
	Quantity        int
	ReserveQuantity int
	Status          ProductStatus
	SoldCount       int
	RatingAvg       float64
	TotalReviews    int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
