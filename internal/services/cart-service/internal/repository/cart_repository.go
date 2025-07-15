package repository

import (
	"gorm.io/gorm"
	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/domain"
)

type CartRepository interface {
	GetCartByUserID(userID uuid.UUID) (*domain.Cart, error)
	CreateCart(cart *domain.Cart) error
	DeleteCart(cartID uuid.UUID) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetCartByUserID(userID uuid.UUID) (*domain.Cart, error) {
	var cart domain.Cart
	if err := r.db.Preload("Items").First(&cart, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) CreateCart(cart *domain.Cart) error {
	return r.db.Create(cart).Error
}

func (r *cartRepository) DeleteCart(cartID uuid.UUID) error {
	return r.db.Delete(&domain.Cart{}, cartID).Error
}