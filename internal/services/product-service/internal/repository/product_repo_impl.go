package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	// 1. Chuyển đổi ID string sang kiểu của pgtype
	productUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid format for repository: %w", err)
	}
	pgUUID := converter.UUIDToPgUUID(productUUID)

	// 2. Gọi query từ SQLC
	sqlcProduct, err := r.queries.GetProductByID(ctx, pgUUID)
	if err != nil {
		// Xử lý trường hợp không tìm thấy record
		if errors.Is(err, pgx.ErrNoRows) {
			// Trả về (nil, nil) là một cách để báo hiệu "không tìm thấy"
			// mà không gây ra lỗi. Application service sẽ diễn giải điều này.
			return nil, nil
		}
		// Đối với các lỗi khác (lỗi kết nối, v.v.), trả về lỗi
		return nil, fmt.Errorf("database error when getting product by id: %w", err)
	}

	// 3. Chuyển đổi từ DB model (sqlc) sang Domain model
	domainProduct, err := toDomain(&sqlcProduct)
	if err != nil {
		return nil, fmt.Errorf("failed to convert sqlc.Product to domain.Product: %w", err)
	}

	return domainProduct, nil
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
