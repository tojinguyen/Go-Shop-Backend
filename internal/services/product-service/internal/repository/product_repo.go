package repository

import (
	"context"

	"github.com/google/uuid"
	product "github.com/toji-dev/go-shop/internal/services/product-service/internal/domain/product"
)

type ProductRepository interface {
	Save(ctx context.Context, product *product.Product) error
	GetByID(ctx context.Context, id string) (*product.Product, error)
	GetByShopID(ctx context.Context, shopID uuid.UUID, limit, offset int) ([]*product.Product, int64, error)
	Update(ctx context.Context, product *product.Product) error
}
