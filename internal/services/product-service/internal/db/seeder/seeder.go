package seeder

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

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

func (s *Seeder) SeedProducts(count int) {
	shopIDs, err := s.fetchShopIDs()
	if err != nil {
		log.Fatalf("❌ Could not fetch shops: %v", err)
	}
	if len(shopIDs) == 0 {
		log.Println("⚠️ No shops found in shop-service DB. Please seed shops first.")
		return
	}

	log.Printf("🌱 Seeding %d products with multi-goroutine COPY FROM...", count)

	// ----- Pre-generate sample data -----
	const sampleSize = 1000
	preNames := make([]string, sampleSize)
	preDescs := make([]string, sampleSize)
	for i := 0; i < sampleSize; i++ {
		preNames[i] = faker.Sentence()
		preDescs[i] = faker.Paragraph()
	}

	// Rand pool để tránh global lock
	var rndPool = sync.Pool{
		New: func() interface{} {
			return rand.New(rand.NewSource(time.Now().UnixNano()))
		},
	}

	// Cấu hình batch và worker
	const batchSize = 5000
	numWorkers := 4 // Nên ≤ pool size của DB

	jobs := make(chan [2]int, numWorkers)
	var wg sync.WaitGroup
	var totalCreated int64
	var mu sync.Mutex

	// Worker goroutine
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rnd := rndPool.Get().(*rand.Rand)
			defer rndPool.Put(rnd)

			for rng := range jobs {
				start, end := rng[0], rng[1]

				rows := make([][]interface{}, 0, batchSize)
				for j := start; j < end; j++ {
					stateChance := rnd.Intn(100)
					var qty int32
					var status sqlc.ProductStatus

					switch {
					case stateChance < 80:
						qty = int32(rnd.Intn(1000) + 10)
						status = sqlc.ProductStatusACTIVE
					case stateChance < 95:
						qty = int32(rnd.Intn(500))
						status = sqlc.ProductStatusINACTIVE
					default:
						qty = 0
						status = sqlc.ProductStatusOUTOFSTOCK
					}

					price := rnd.Intn(4990001) + 10000
					shopID := shopIDs[rnd.Intn(len(shopIDs))]

					rows = append(rows, []interface{}{
						shopID,
						preNames[rnd.Intn(sampleSize)],
						preDescs[rnd.Intn(sampleSize)],
						float64(price),
						"VND",
						qty,
						0,
						status,
					})
				}

				// COPY FROM batch
				copyCount, err := s.productDB.CopyFrom(
					s.ctx,
					pgx.Identifier{"products"},
					[]string{
						"shop_id", "product_name", "product_description",
						"price", "currency", "quantity", "reserve_quantity", "product_status",
					},
					pgx.CopyFromRows(rows),
				)
				if err != nil {
					log.Printf("❌ Error batch %d-%d: %v", start+1, end, err)
					continue
				}

				mu.Lock()
				totalCreated += int64(copyCount)
				mu.Unlock()

				log.Printf("✅ Batch %d-%d done. Total: %d/%d", start+1, end, totalCreated, count)
			}
		}()
	}

	// Gửi job vào channel
	for i := 0; i < count; i += batchSize {
		end := i + batchSize
		if end > count {
			end = count
		}
		jobs <- [2]int{i, end}
	}
	close(jobs)

	wg.Wait()
	log.Printf("🎉 Done seeding. Total: %d products", totalCreated)
}
