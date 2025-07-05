package repository

import (
	"context"

	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/domain"
)

// ShopRepository defines the interface for shop data operations
type ShopRepository interface {
	// Create creates a new shop
	Create(ctx context.Context, shop *domain.Shop) error
}
