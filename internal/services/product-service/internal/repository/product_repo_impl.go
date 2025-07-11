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

func (r *pgProductRepository) GetByShopID(ctx context.Context, shopID uuid.UUID, limit, offset int) ([]*product.Product, int64, error) {
	pgShopUUID := converter.UUIDToPgUUID(shopID)

	// Sử dụng transaction để đảm bảo 2 câu query là nhất quán
	tx, err := r.db.BeginTransaction(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Rollback nếu có lỗi xảy ra

	qtx := r.queries.WithTx(tx)

	// 1. Lấy tổng số sản phẩm
	totalCount, err := qtx.CountProductsByShop(ctx, pgShopUUID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products by shop: %w", err)
	}

	if totalCount == 0 {
		return []*product.Product{}, 0, nil
	}

	// 2. Lấy danh sách sản phẩm theo phân trang
	params := sqlc.GetListProductsByShopParams{
		ShopID: pgShopUUID,
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	sqlcProducts, err := qtx.GetListProductsByShop(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get list of products by shop: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 3. Chuyển đổi kết quả sang domain models
	domainProducts := make([]*product.Product, 0, len(sqlcProducts))
	for _, p := range sqlcProducts {
		domainProduct, err := toDomain(&p)
		if err != nil {
			// Có thể log lỗi và bỏ qua sản phẩm này hoặc trả về lỗi ngay lập tức
			// Tạm thời trả về lỗi để đảm bảo tính toàn vẹn
			return nil, 0, fmt.Errorf("failed to convert product %s: %w", p.ID, err)
		}
		domainProducts = append(domainProducts, domainProduct)
	}

	return domainProducts, totalCount, nil
}

func (r *pgProductRepository) Update(ctx context.Context, p *product.Product) error {
	params := sqlc.UpdateProductParams{
		ID:                 converter.UUIDToPgUUID(p.ID()),
		ProductName:        p.Name(),
		ProductDescription: converter.StringToPgText(p.Description()),
		CategoryID:         converter.UUIDToPgUUID(p.CategoryID()),
		Price:              converter.Float64ToPgNumeric(p.Price().GetAmount()),
		Currency:           p.Price().GetCurrency(),
		Quantity:           int32(p.Quantity()),
		ThumbnailUrl:       converter.StringToPgText(p.ThumbnailURL()),
		ProductStatus:      sqlc.ProductStatus(p.Status()),
	}

	_, err := r.queries.UpdateProduct(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update product in db: %w", err)
	}

	return nil
}

func (r *pgProductRepository) Delete(ctx context.Context, id string) error {
	productUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid uuid format for repository: %w", err)
	}
	pgUUID := converter.UUIDToPgUUID(productUUID)

	// Gọi hàm SoftDeleteProduct đã được sqlc generate
	err = r.queries.SoftDeleteProduct(ctx, pgUUID)
	if err != nil {
		return fmt.Errorf("failed to soft delete product in db: %w", err)
	}

	return nil
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
