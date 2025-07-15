package repository

import (
	"gorm.io/gorm"
	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/domain"
)

type CartItemRepository interface {
	AddItemToCart(item *domain.CartItem) error
	UpdateCartItem(item *domain.CartItem) error
	RemoveCartItem(cartID, productID uuid.UUID) error
}

type cartItemRepository struct {
	db *gorm.DB
}

func NewCartItemRepository(db *gorm.DB) CartItemRepository {
	return &cartItemRepository{db: db}
}

func (r *cartItemRepository) AddItemToCart(item *domain.CartItem) error {
	return r.db.Create(item).Error
}

func (r *cartItemRepository) UpdateCartItem(item *domain.CartItem) error {
	return r.db.Model(&domain.CartItem{}).Where("cart_id = ? AND product_id = ?", item.CartID, item.ProductID).Updates(item).Error
}

func (r *cartItemRepository) RemoveCartItem(cartID, productID uuid.UUID) error {
	return r.db.Delete(&domain.CartItem{}, "cart_id = ? AND product_id = ?", cartID, productID).Error
}