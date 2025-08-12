package seeder

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/db/sqlc"
)

type Seeder struct {
	productDB *pgxpool.Pool
	shopDB    *pgxpool.Pool
	ctx       context.Context
}

func NewSeeder(productDB, shopDB *pgxpool.Pool) *Seeder {
	return &Seeder{
		productDB: productDB,
		shopDB:    shopDB,
		ctx:       context.Background(),
	}
}

// fetchShopIDs vẫn giữ nguyên
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

// SeedProducts đã được tối ưu hóa với pgx.CopyFrom
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

	log.Printf("🌱 Seeding %d products using highly optimized 'COPY' protocol...", count)

	// [TỐI ƯU HÓA] Định nghĩa tên các cột sẽ được chèn.
	// Thứ tự phải khớp với thứ tự các giá trị trong mỗi row.
	columnNames := []string{
		"shop_id",
		"product_name",
		"thumbnail_url",
		"product_description",
		"category_id",
		"price",
		"currency",
		"quantity",
		"reserve_quantity",
		"product_status",
	}

	const batchSize = 1000 // Có thể tăng lên 5000 hoặc 10000 để nhanh hơn nữa
	productsCreated := 0

	for i := 0; i < count; i += batchSize {
		batchEnd := i + batchSize
		if batchEnd > count {
			batchEnd = count
		}

		log.Printf("Preparing batch %d-%d...", i+1, batchEnd)

		// [TỐI ƯU HÓA] Tạo một slice chứa các hàng dữ liệu cho batch này.
		rows := make([][]interface{}, 0, batchSize)

		for j := i; j < batchEnd; j++ {
			var quantity int32
			var status sqlc.ProductStatus
			stateChance := rand.Intn(100)

			switch {
			case stateChance < 80:
				quantity = int32(rand.Intn(1000) + 10)
				status = sqlc.ProductStatusACTIVE
			case stateChance < 95:
				quantity = int32(rand.Intn(500))
				status = sqlc.ProductStatusINACTIVE
			default:
				quantity = 0
				status = sqlc.ProductStatusOUTOFSTOCK
			}

			shopID := shopIDs[rand.Intn(len(shopIDs))]
			productDesc := faker.Paragraph()
			price, _ := faker.RandomInt(10000, 5000000, 1)

			// [TỐI ƯU HÓA] Thêm một hàng dữ liệu vào slice.
			// Lưu ý: Thứ tự phải khớp với `columnNames` đã định nghĩa ở trên.
			rows = append(rows, []interface{}{
				shopID,
				faker.Sentence(),
				nil, // thumbnail_url
				productDesc,
				nil, // category_id
				float64(price[0]),
				"VND",
				quantity,
				0, // reserve_quantity
				status,
			})
		}

		// [TỐI ƯU HÓA] Sử dụng CopyFrom để chèn toàn bộ batch vào DB.
		copyCount, err := s.productDB.CopyFrom(
			s.ctx,
			pgx.Identifier{"products"},
			columnNames,
			pgx.CopyFromRows(rows),
		)

		if err != nil {
			log.Printf("❌ Error processing batch %d-%d: %v. Skipping this batch.", i+1, batchEnd, err)
			continue
		}

		if int(copyCount) != len(rows) {
			log.Printf("⚠️ Mismatch count for batch %d-%d: expected %d, got %d. Some rows might not have been inserted.", i+1, batchEnd, len(rows), copyCount)
		}

		productsCreated += int(copyCount)
		log.Printf("✅ Successfully seeded batch %d-%d. Total seeded: %d/%d", i+1, batchEnd, productsCreated, count)
	}

	log.Printf("🎉 Product seeding complete. Total products created: %d", productsCreated)
}
