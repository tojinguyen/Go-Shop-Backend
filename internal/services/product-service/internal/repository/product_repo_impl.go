package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/pkg/converter"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/db/sqlc"
	product "github.com/toji-dev/go-shop/internal/services/product-service/internal/domain/product"
	product_v1 "github.com/toji-dev/go-shop/proto/gen/go/product/v1"
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

func (r *pgProductRepository) Save(ctx context.Context, p *product.Product) (*product.Product, error) {
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

	product, err := r.queries.CreateProduct(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to save product to db: %w", err)
	}

	// Chuyển đổi từ DB model (sqlc) sang Domain model
	productDomain, err := toDomain(&product)
	if err != nil {
		return nil, fmt.Errorf("failed to convert db model to domain: %w", err)
	}

	return productDomain, nil
}

func (r *pgProductRepository) GetByID(ctx context.Context, id string) (*product.Product, error) {
	// 1. Chuyển đổi ID string sang kiểu của pgtype
	productUUID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Error parsing product ID: %v", err)
		return nil, fmt.Errorf("invalid uuid format for repository: %w", err)
	}
	pgUUID := converter.UUIDToPgUUID(productUUID)

	// 2. Gọi query từ SQLC
	sqlcProduct, err := r.queries.GetProductByID(ctx, pgUUID)
	if err != nil {
		log.Printf("Error getting product by ID from db: %v", err)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("Product with ID %s not found", id)
			return nil, apperror.NewNotFound("product", id)
		}
		log.Printf("Database error when getting product by ID: %v", err)
		return nil, apperror.New(apperror.CodeDatabaseError, "database error when getting product by id", apperror.TypeInternal).Wrap(err)
	}

	// 3. Chuyển đổi từ DB model (sqlc) sang Domain model
	domainProduct, err := toDomain(&sqlcProduct)
	if err != nil {
		log.Printf("Error converting SQLC product to domain product: %v", err)
		return nil, apperror.New(apperror.CodeConversionError, "failed to convert db model to domain", apperror.TypeInternal).Wrap(err)
	}

	return domainProduct, nil
}

func (r *pgProductRepository) GetByShopID(ctx context.Context, shopID uuid.UUID, limit, offset int) ([]*product.Product, int64, error) {
	pgShopUUID := converter.UUIDToPgUUID(shopID)

	// Sử dụng transaction để đảm bảo 2 câu query là nhất quán
	tx, err := r.db.BeginTransaction(ctx)
	if err != nil {
		return nil, 0, apperror.New(apperror.CodeDatabaseError, "failed to begin transaction", apperror.TypeInternal).Wrap(err)
	}
	defer tx.Rollback(ctx) // Rollback nếu có lỗi xảy ra

	qtx := r.queries.WithTx(tx)

	// 1. Lấy tổng số sản phẩm
	totalCount, err := qtx.CountProductsByShop(ctx, pgShopUUID)
	if err != nil {
		return nil, 0, apperror.New(apperror.CodeDatabaseError, "failed to count products by shop", apperror.TypeInternal).Wrap(err)
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
		return nil, 0, apperror.New(apperror.CodeDatabaseError, "failed to get products by shop", apperror.TypeInternal).Wrap(err)
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
			return nil, 0, apperror.New(apperror.CodeConversionError, "failed to convert db model to domain", apperror.TypeInternal).Wrap(err)
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
		return apperror.New(apperror.CodeDatabaseError, "failed to update product in db", apperror.TypeInternal).Wrap(err)
	}

	return nil
}

func (r *pgProductRepository) Delete(ctx context.Context, id string) error {
	productUUID, err := uuid.Parse(id)
	if err != nil {
		return apperror.New(apperror.CodeConversionError, "invalid uuid format for repository", apperror.TypeInternal).Wrap(err)
	}
	pgUUID := converter.UUIDToPgUUID(productUUID)

	// Gọi hàm SoftDeleteProduct đã được sqlc generate
	err = r.queries.SoftDeleteProduct(ctx, pgUUID)
	if err != nil {
		return apperror.New(apperror.CodeDatabaseError, "failed to delete product from db", apperror.TypeInternal).Wrap(err)
	}

	return nil
}

func toDomain(p *sqlc.Product) (*product.Product, error) {
	price, err := product.NewPrice(
		*converter.PgNumericToFloat64Ptr(p.Price),
		p.Currency,
	)
	if err != nil {
		return nil, apperror.New(apperror.CodeConversionError, "failed to convert price", apperror.TypeInternal).Wrap(err)
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

func (r *pgProductRepository) GetByIDs(ctx context.Context, ids []string) ([]*product.Product, error) {
	if len(ids) == 0 {
		return []*product.Product{}, nil
	}

	pgIDs := make([]pgtype.UUID, len(ids))
	for i, id := range ids {
		pgIDs[i] = converter.UUIDToPgUUID(uuid.Must(uuid.Parse(id)))
	}

	sqlcProducts, err := r.queries.GetProductsByIDs(ctx, pgIDs)
	if err != nil {
		return nil, apperror.New(apperror.CodeDatabaseError, "failed to get products by IDs", apperror.TypeInternal).Wrap(err)
	}

	domainProducts := make([]*product.Product, 0, len(sqlcProducts))
	for _, p := range sqlcProducts {
		domainProduct, err := toDomain(&p)
		if err != nil {
			return nil, apperror.New(apperror.CodeConversionError, "failed to convert db model to domain", apperror.TypeInternal).Wrap(err)
		}
		domainProducts = append(domainProducts, domainProduct)
	}

	return domainProducts, nil
}

func (r *pgProductRepository) ReserveStock(ctx context.Context, items []*product_v1.ReserveProduct) ([]*product_v1.ProductReservationStatus, error) {
	// Bắt đầu một giao dịch CSDL.
	tx, err := r.db.BeginTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Chuẩn bị dữ liệu để query và kiểm tra
	productIDs := make([]pgtype.UUID, len(items))
	requestedQuantityMap := make(map[string]int32, len(items))
	for i, item := range items {
		productIDs[i] = converter.StringToPgUUID(item.ProductId)
		requestedQuantityMap[item.ProductId] = item.Quantity
	}

	// BƯỚC 1: LẤY VÀ KHOÁ CÁC SẢN PHẨM
	dbProducts, err := qtx.GetProductsByIDsForUpdate(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get and lock products: %w", err)
	}

	// Chuyển đổi sang map để dễ dàng truy cập
	productMap := make(map[string]sqlc.Product, len(dbProducts))
	for _, p := range dbProducts {
		productMap[p.ID.String()] = p
	}

	statuses := make([]*product_v1.ProductReservationStatus, len(items))
	allSuccess := true

	// BƯỚC 2: KIỂM TRA TỒN KHO TRONG GIAO DỊCH
	// Vòng lặp này chỉ để kiểm tra và tạo ra các thông báo trạng thái.
	for i, item := range items {
		product, ok := productMap[item.ProductId]
		if !ok {
			allSuccess = false
			statuses[i] = &product_v1.ProductReservationStatus{
				ProductId: item.ProductId,
				Success:   false,
				Message:   "Product not found",
			}
			continue
		}

		availableStock := product.Quantity - product.ReserveQuantity
		if availableStock < item.Quantity {
			allSuccess = false
			statuses[i] = &product_v1.ProductReservationStatus{
				ProductId: item.ProductId,
				Success:   false,
				Message:   fmt.Sprintf("Insufficient stock. Available: %d, Requested: %d", availableStock, item.Quantity),
			}
			continue
		}

		// Nếu đủ hàng, tạo trạng thái thành công
		statuses[i] = &product_v1.ProductReservationStatus{
			ProductId: item.ProductId,
			Success:   true,
			Message:   "Reserved",
		}
	}

	// BƯỚC 3: QUYẾT ĐỊNH COMMIT HAY ROLLBACK
	if !allSuccess {
		log.Printf("Stock reservation failed for one or more items. Rolling back transaction.")
		return statuses, nil // Trả về lỗi nghiệp vụ, không phải lỗi hệ thống
	}

	// BƯỚC 4: THỰC HIỆN CẬP NHẬT
	for _, item := range items {
		product := productMap[item.ProductId]
		newReserveQuantity := product.ReserveQuantity + item.Quantity

		_, err := qtx.UpdateProductStock(ctx, sqlc.UpdateProductStockParams{
			ID:              product.ID,
			Quantity:        product.Quantity,
			ReserveQuantity: newReserveQuantity,
		})

		if err != nil {
			return nil, fmt.Errorf("failed to update stock for product %s: %w", item.ProductId, err)
		}
	}

	// BƯỚC 5: COMMIT GIAO DỊCH
	// Nếu mọi thứ thành công, commit để lưu lại tất cả các thay đổi.
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit stock reservation transaction: %w", err)
	}

	log.Println("Stock reservation transaction committed successfully.")
	return statuses, nil
}

func (r *pgProductRepository) IsOrderReserved(ctx context.Context, orderID string) (bool, error) {
	pgOrderID := converter.StringToPgUUID(orderID)

	// Sử dụng query để kiểm tra xem đơn hàng đã được đặt trước hay chưa
	isReserved, err := r.queries.IsOrderReserved(ctx, pgOrderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check order reservation: %w", err)
	}

	return isReserved, nil
}
