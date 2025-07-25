package repository

import (
	"context"

	"github.com/google/uuid"
	product "github.com/toji-dev/go-shop/internal/services/product-service/internal/domain/product"
	product_v1 "github.com/toji-dev/go-shop/proto/gen/go/product/v1"
)

type ProductRepository interface {
	Save(ctx context.Context, product *product.Product) (*product.Product, error)
	GetByID(ctx context.Context, id string) (*product.Product, error)
	GetByShopID(ctx context.Context, shopID uuid.UUID, limit, offset int) ([]*product.Product, int64, error)
	Update(ctx context.Context, product *product.Product) error
	Delete(ctx context.Context, id string) error
	GetByIDs(ctx context.Context, ids []string) ([]*product.Product, error)
	ReserveStock(ctx context.Context, items []*product_v1.ReserveProduct) ([]*product_v1.ProductReservationStatus, error)
	IsOrderReserved(ctx context.Context, orderID string) (bool, error)
}
