package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Price is a Value Object
type Price struct {
	Amount   float64
	Currency string
}

// Product is the Aggregate Root
type Product struct {
	ID            uuid.UUID
	ShopID        uuid.UUID
	Name          string
	Description   string
	Price         Price
	StockQuantity int
	CategoryID    uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewProduct is a factory function for creating a new product
func NewProduct(shopID, categoryID uuid.UUID, name, description, currency string, price float64, stock int) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	if price <= 0 {
		return nil, errors.New("product price must be positive")
	}

	return &Product{
		ID:            uuid.New(),
		ShopID:        shopID,
		Name:          name,
		Description:   description,
		Price:         Price{Amount: price, Currency: currency},
		StockQuantity: stock,
		CategoryID:    categoryID,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}, nil
}

func (p *Product) ChangePrice(newPrice float64) error {
	if newPrice <= 0 {
		return errors.New("price must be positive")
	}
	p.Price.Amount = newPrice
	p.UpdatedAt = time.Now().UTC()
	return nil
}
