package domain

import (
	"errors"
	"time"

	time_utils "github.com/toji-dev/go-shop/internal/pkg/time"

	uuid "github.com/google/uuid"
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

func NewProduct(shopID, name, thumbnailURL, description string, categoryID *uuid.UUID, price Price, quantity int) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	return &Product{
		ID:              uuid.New(),
		ShopID:          uuid.MustParse(shopID),
		Name:            name,
		ThumbnailURL:    thumbnailURL,
		Description:     description,
		CategoryID:      categoryID,
		Price:           price,
		Quantity:        quantity,
		ReserveQuantity: 0,
		Status:          ProductStatusDraft,
		SoldCount:       0,
		RatingAvg:       0.0,
		TotalReviews:    0,
		CreatedAt:       time_utils.GetUtcTime(),
		UpdatedAt:       time_utils.GetUtcTime(),
	}, nil
}

func (p *Product) ChangePrice(newPrice Price) error {
	if newPrice.Amount < 0 {
		return errors.New("price cannot be negative")
	}

	p.Price = newPrice
	p.UpdatedAt = time_utils.GetUtcTime()
	return nil
}

func (p *Product) UpdateQuantity(newQuantity int) error {
	if newQuantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	p.Quantity = newQuantity
	p.UpdatedAt = time_utils.GetUtcTime()
	return nil
}

func (p *Product) Deactivate() {
	p.Status = ProductStatusInactive
	p.UpdatedAt = time_utils.GetUtcTime()
}
