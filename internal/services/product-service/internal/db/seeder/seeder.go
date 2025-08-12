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

// SeedProducts đã được tối ưu hóa với pre-generation và pgx.CopyFrom
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

	log.Printf("🌱 Seeding %d products using highly optimized 'COPY' protocol and pre-generation...", count)

	// [TỐI ƯU HÓA] Bước 1: Tạo sẵn một bộ dữ liệu mẫu để tránh gọi faker trong vòng lặp lớn
	log.Println("Pre-generating sample data...")
	const sampleSize = 200 // Tạo 200 mẫu tên và mô tả
	preGeneratedNames := make([]string, sampleSize)
	preGeneratedDescriptions := make([]string, sampleSize)
	for i := 0; i < sampleSize; i++ {
		preGeneratedNames[i] = faker.Sentence()
		preGeneratedDescriptions[i] = faker.Paragraph()
	}
	log.Println("Sample data generated.")

	// Định nghĩa tên các cột sẽ được chèn.
	columnNames := []string{
		"shop_id",
		"product_name",
		"product_description",
		"price",
		"currency",
		"quantity",
		"reserve_quantity",
		"product_status",
	}

	const batchSize = 1000 // Tăng batch size để hiệu quả hơn
	productsCreated := 0

	for i := 0; i < count; i += batchSize {
		batchEnd := i + batchSize
		if batchEnd > count {
			batchEnd = count
		}

		rows := make([][]interface{}, 0, batchSize)

		// [TỐI ƯU HÓA] Bước 2: Tạo dữ liệu cho batch từ các mẫu đã có, cực kỳ nhanh
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
			// Sử dụng math/rand thay vì faker.RandomInt để nhanh hơn
			price := rand.Intn(4990001) + 10000 // Giá từ 10,000 đến 5,000,000

			// Lấy dữ liệu từ bộ nhớ thay vì tạo mới
			productName := preGeneratedNames[rand.Intn(sampleSize)]
			productDesc := preGeneratedDescriptions[rand.Intn(sampleSize)]

			rows = append(rows, []interface{}{
				shopID,
				productName,
				productDesc,
				float64(price),
				"VND",
				quantity,
				0, // reserve_quantity
				status,
			})
		}

		// [TỐI ƯU HÓA] Bước 3: Sử dụng CopyFrom để chèn toàn bộ batch
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
			log.Printf("⚠️ Mismatch count for batch %d-%d: expected %d, got %d.", i+1, batchEnd, len(rows), copyCount)
		}

		productsCreated += int(copyCount)
		log.Printf("✅ Successfully seeded batch %d-%d. Total seeded: %d/%d", i+1, batchEnd, productsCreated, count)
	}

	log.Printf("🎉 Product seeding complete. Total products created: %d", productsCreated)
}
