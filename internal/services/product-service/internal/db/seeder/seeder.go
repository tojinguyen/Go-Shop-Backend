package seeder

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/toji-dev/go-shop/internal/pkg/converter"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/db/sqlc"
)

type Seeder struct {
	productDB *pgxpool.Pool
	shopDB    *pgxpool.Pool
	queries   *sqlc.Queries
	ctx       context.Context
}

func NewSeeder(productDB, shopDB *pgxpool.Pool) *Seeder {
	return &Seeder{
		productDB: productDB,
		shopDB:    shopDB,
		queries:   sqlc.New(productDB),
		ctx:       context.Background(),
	}
}

// fetchShopIDs lấy danh sách ID của các shop từ DB của shop-service
func (s *Seeder) fetchShopIDs() ([]uuid.UUID, error) {
	log.Println("Fetching shop IDs from shop-service database...")
	rows, err := s.shopDB.Query(s.ctx, "SELECT id FROM shops")
	if err != nil {
		return nil, fmt.Errorf("failed to query shop IDs: %w", err)
	}
	defer rows.Close()

	var shopIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan shop ID: %w", err)
		}
		shopIDs = append(shopIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over shop IDs: %w", err)
	}

	log.Printf("Found %d shops.", len(shopIDs))
	return shopIDs, nil
}

// SeedProducts tạo dữ liệu giả cho sản phẩm với logic đã được cải tiến
func (s *Seeder) SeedProducts(count int) {
	shopIDs, err := s.fetchShopIDs()
	if err != nil {
		log.Fatalf("❌ Could not fetch shops: %v", err)
	}
	if len(shopIDs) == 0 {
		log.Println("⚠️ No shops found in shop-service DB. Please seed shops first.")
		log.Println("Run: make seed-shops")
		return
	}

	log.Printf("🌱 Seeding %d products with improved logic...", count)

	// [XÓA BỎ] Mảng trạng thái ngẫu nhiên đã được loại bỏ.
	// productStatuses := []sqlc.ProductStatus{ ... }

	for i := 0; i < count; i++ {
		// Chọn ngẫu nhiên một shop
		shopID := shopIDs[rand.Intn(len(shopIDs))]

		// =======================================================
		// [LOGIC MỚI] Xác định trạng thái và số lượng sản phẩm một cách logic
		// =======================================================
		var quantity int32
		var status sqlc.ProductStatus

		// Phân phối trạng thái sản phẩm để dữ liệu thực tế hơn
		stateChance := rand.Intn(100) // Tạo số ngẫu nhiên từ 0-99

		switch {
		case stateChance < 80: // 80% trường hợp: Sản phẩm đang hoạt động và có hàng
			quantity = int32(rand.Intn(1000) + 10) // Số lượng tồn kho từ 10 đến 1009
			status = sqlc.ProductStatusACTIVE
		case stateChance < 95: // 15% trường hợp: Sản phẩm không hoạt động (người bán tạm ẩn)
			quantity = int32(rand.Intn(500)) // Có thể có hoặc không có hàng
			status = sqlc.ProductStatusINACTIVE
		default: // 5% trường hợp còn lại: Hết hàng
			quantity = 0
			status = sqlc.ProductStatusOUTOFSTOCK
		}
		// =======================================================

		// Tạo dữ liệu sản phẩm giả
		productDesc := faker.Paragraph()
		price, _ := faker.RandomInt(10000, 5000000, 1) // Giá từ 10k đến 5tr

		params := sqlc.CreateProductParams{
			ShopID:             converter.UUIDToPgUUID(shopID),
			ProductName:        faker.Sentence(),
			ThumbnailUrl:       converter.StringToPgText(nil), // Có thể thêm URL ảnh giả ở đây
			ProductDescription: converter.StringToPgText(&productDesc),
			Price:              converter.Float64ToPgNumeric(float64(price[0])),
			Currency:           "VND",
			Quantity:           quantity, // [THAY ĐỔI] Sử dụng số lượng đã được quyết định ở trên
			ReserveQuantity:    0,        // Sản phẩm mới tạo chưa có ai đặt trước
			ProductStatus:      status,   // [THAY ĐỔI] Sử dụng trạng thái logic
		}

		_, err := s.queries.CreateProduct(s.ctx, params)
		if err != nil {
			log.Printf("Failed to create product for shop %s: %v", shopID, err)
			continue // Bỏ qua sản phẩm này và tiếp tục
		}

		if (i+1)%100 == 0 {
			log.Printf("... seeded %d/%d products", i+1, count)
		}
	}

	log.Println("🎉 Product seeding complete.")
}
