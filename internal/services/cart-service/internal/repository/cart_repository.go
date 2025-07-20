package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/domain"
	"gorm.io/gorm"
)

type CartRepository interface {
	GetCartByUserID(ctx context.Context, userID uuid.UUID) (*domain.Cart, error)
	CreateCart(ctx context.Context, cart *domain.Cart) error
	DeleteCart(ctx context.Context, cartID uuid.UUID) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetCartByUserID(ctx context.Context, userID uuid.UUID) (*domain.Cart, error) {
	var cart domain.Cart
	if err := r.db.Preload("Items").First(&cart, "user_id = ?", userID).Error; err != nil {
		log.Printf("Error fetching cart for user %s: %v", userID, err)
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.NewNotFound("cart", userID.String())
		}
		return nil, apperror.NewInternal(string(apperror.CodeDatabaseError))
	}
	return &cart, nil
}

func (r *cartRepository) CreateCart(ctx context.Context, cart *domain.Cart) error {
	return r.db.Create(cart).Error
}

func (r *cartRepository) DeleteCart(ctx context.Context, cartID uuid.UUID) error {
	return r.db.Delete(&domain.Cart{}, cartID).Error
}
