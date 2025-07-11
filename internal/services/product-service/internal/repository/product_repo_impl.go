package repository

import (
	"context"
	"fmt"

	"github.com/toji-dev/go-shop/internal/pkg/converter"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/db/sqlc"
	domain "github.com/toji-dev/go-shop/internal/services/product-service/internal/domain/product"
)

type pgProductRepository struct {
	db      *postgresql_infra.PostgreSQLService
	queries *sqlc.Queries
}

func NewProductRepository(db *postgresql_infra.PostgreSQLService) ProductRepository {
	if db == nil {
		return nil
	}

	queries := sqlc.New(db.GetPool())
	return &pgProductRepository{
		db:      db,
		queries: queries,
	}
}

func (r *pgProductRepository) Save(ctx context.Context, product *domain.Product) error {
	params := sqlc.CreateProductParams{
		ShopID:             converter.UUIDToPgUUID(product.ShopID()),
		ProductName:        product.Name(),
		ThumbnailUrl:       converter.StringToPgText(product.ThumbnailURL()),
		ProductDescription: converter.StringToPgText(product.Description()),
		CategoryID:         converter.UUIDToPgUUID(product.CategoryID()),
		Price:              converter.Float64ToPgNumeric(product.Price().GetAmount()),
		Quantity:           int32(product.Quantity()),
		ReserveQuantity:    int32(product.Quantity()),
		ProductStatus:      sqlc.ProductStatus(product.Status()),
	}

	_, err := r.queries.CreateProduct(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to save product: %w", err)
	}

	return nil
}

func (r *pgProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	product, err := r.queries.GetProductByID(ctx, converter.StringToPgUUID(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}

	price, err := domain.NewPrice(
		*converter.PgNumericToFloat64Ptr(product.Price),
		"USD",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create price from product data: %w", err)
	}

	productDomain, err := domain.NewProduct(
		product.ShopID.String(),
		product.ProductName,
		product.ThumbnailUrl.String,
		product.ProductDescription.String,
		converter.PgUUIDToUUID(product.CategoryID),
		price,
		int(product.Quantity),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to convert product to domain: %w", err)
	}

	return productDomain, nil
}

func (r *pgProductRepository) GetByShopID(ctx context.Context, shopID string) ([]*domain.Product, error) {
	return nil, nil
}
