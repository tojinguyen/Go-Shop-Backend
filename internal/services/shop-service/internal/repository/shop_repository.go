package repository

import (
	"context"

	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/domain"
)

// ShopRepository defines the interface for shop data operations
type ShopRepository interface {
	Create(ctx context.Context, shop *domain.Shop) error
	GetShopByID(ctx context.Context, shopID string) (*domain.Shop, error)
	GetShopsByOwnerID(ctx context.Context, ownerID string) ([]*domain.Shop, error)
	Update(ctx context.Context, shop *domain.Shop) error
	Delete(ctx context.Context, shopID string) error
}
