package repository

import (
	"context"

	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/domain"
)

type PromotionRepository interface {
	Create(ctx context.Context, promotion *domain.Promotion) error
	GetByID(ctx context.Context, id string) (*domain.Promotion, error)
	GetByShopID(ctx context.Context, shopID string) ([]*domain.Promotion, error)
	Update(ctx context.Context, promotion *domain.Promotion) error
	Delete(ctx context.Context, id string) error
}
