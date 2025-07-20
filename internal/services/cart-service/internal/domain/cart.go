package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	time_utils "github.com/toji-dev/go-shop/internal/pkg/time"
)

type Cart struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Items     []CartItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewCart là factory function
func NewCart(ownerID uuid.UUID) *Cart {
	return &Cart{
		ID:     uuid.New(),
		UserID: ownerID,
		Items:  []CartItem{},
	}
}

// findItem tìm một item trong giỏ hàng bằng productID
func (c *Cart) findItem(productID uuid.UUID) *CartItem {
	for i := range c.Items {
		if c.Items[i].ProductID == productID {
			return &c.Items[i]
		}
	}
	return nil
}

// AddItem chứa logic nghiệp vụ thêm sản phẩm
func (c *Cart) AddItem(productID uuid.UUID, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	// Logic bạn muốn refactor giờ nằm ở đây
	existingItem := c.findItem(productID)

	if existingItem != nil {
		if err := existingItem.IncreaseQuantity(quantity); err != nil {
			return err
		}
	} else {
		newItem, err := NewCartItem(c.ID, productID, quantity)
		if err != nil {
			return err
		}
		c.Items = append(c.Items, *newItem)
	}

	c.UpdatedAt = time_utils.GetUtcTime()
	return nil
}
