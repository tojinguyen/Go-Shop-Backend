package repository

import (
	"context"
	"fmt"

	"github.com/toji-dev/go-shop/internal/pkg/converter"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/db/sqlc"
	product "github.com/toji-dev/go-shop/internal/services/product-service/internal/domain/product"
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

func (r *pgProductRepository) Save(ctx context.Context, p *product.Product) error {
	params := sqlc.CreateProductParams{
		ShopID:             converter.UUIDToPgUUID(p.ShopID()),
		ProductName:        p.Name(),
		ThumbnailUrl:       converter.StringToPgText(p.ThumbnailURL()),
		ProductDescription: converter.StringToPgText(p.Description()),
		CategoryID:         converter.UUIDToPgUUID(p.CategoryID()),
		Price:              converter.Float64ToPgNumeric(p.Price().GetAmount()),
		Currency:           p.Price().GetCurrency(),
		Quantity:           int32(p.Quantity()),
		ReserveQuantity:    0,
		ProductStatus:      sqlc.ProductStatus(p.Status()),
	}

	_, err := r.queries.CreateProduct(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to save product to db: %w", err)
	}
	return nil
}

func (r *pgProductRepository) GetByID(ctx context.Context, id string) (*product.Product, error) {
	productRes, err := r.queries.GetProductByID(ctx, converter.StringToPgUUID(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}

	productDomain, err := toDomain(&productRes)

	if err != nil {
		return nil, fmt.Errorf("failed to convert product to domain: %w", err)
	}

	return productDomain, nil
}

func (r *pgProductRepository) GetProductsByShopID(ctx context.Context, shopID string) ([]*product.Product, error) {
	return nil, nil
}

func toDomain(p *sqlc.Product) (*product.Product, error) {
	price, err := product.NewPrice(
		*converter.PgNumericToFloat64Ptr(p.Price),
		p.Currency,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid price data from db: %w", err)
	}

	return product.Reconstitute(
		converter.PgUUIDToUUID(p.ID),
		converter.PgUUIDToUUID(p.ShopID),
		converter.PgUUIDToUUID(p.CategoryID),
		p.ProductName,
		p.ProductDescription.String,
		p.ThumbnailUrl.String,
		price,
		int(p.Quantity),
		product.ProductStatus(p.ProductStatus),
		p.CreatedAt.Time,
		p.UpdatedAt.Time,
	)
}
