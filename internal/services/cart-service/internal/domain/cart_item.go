package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CartID    uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	ShopID    uuid.UUID `gorm:"type:uuid;not null"`
	Quantity  int       `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewCartItem là một factory function để tạo CartItem mới
func NewCartItem(cartID, productID uuid.UUID, quantity int) (*CartItem, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}
	return &CartItem{
		ID:        uuid.New(),
		CartID:    cartID,
		ProductID: productID,
		Quantity:  quantity,
	}, nil
}

func (item *CartItem) IncreaseQuantity(amount int) error {
	if amount <= 0 {
		return errors.New("amount to add must be positive")
	}
	item.Quantity += amount
	return nil
}
