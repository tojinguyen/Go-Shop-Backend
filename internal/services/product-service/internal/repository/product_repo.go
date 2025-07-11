package repository

import (
	"context"

	product "github.com/toji-dev/go-shop/internal/services/product-service/internal/domain/product"
)

type ProductRepository interface {
	Save(ctx context.Context, product *product.Product) error
	GetByID(ctx context.Context, id string) (*product.Product, error)
	GetProductsByShopID(ctx context.Context, shopID string) ([]*product.Product, error)
}
