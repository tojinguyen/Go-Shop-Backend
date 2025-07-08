package repository

import (
	"context"

	"github.com/toji-dev/go-shop/internal/services/product-service/internal/domain"
)

type ProductRepository interface {
	Save(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	GetByShopID(ctx context.Context, shopID string) ([]*domain.Product, error)
}
