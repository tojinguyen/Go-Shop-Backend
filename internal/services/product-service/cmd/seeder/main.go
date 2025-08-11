package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/db/seeder"
)

func main() {
	productCount := flag.Int("products", 1000, "Number of products to seed")
	flag.Parse()

	log.Println("🌱 Starting product-service database seeder...")

	// 1. Tải cấu hình của product-service
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// 2. Kết nối đến DB của product-service
	productDB, err := pgxpool.New(context.Background(), cfg.GetDatabaseURL())
	if err != nil {
		log.Fatalf("❌ Failed to connect to product-service database: %v", err)
	}
	defer productDB.Close()

	// 3. Kết nối đến DB của shop-service
	shopDBConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.ShopServiceDB.User,
		cfg.ShopServiceDB.Password,
		cfg.ShopServiceDB.Host,
		cfg.ShopServiceDB.Port,
		cfg.ShopServiceDB.DBName,
		cfg.ShopServiceDB.SSLMode,
	)
	shopDB, err := pgxpool.New(context.Background(), shopDBConnStr)
	if err != nil {
		log.Fatalf("❌ Failed to connect to shop-service database: %v", err)
	}
	defer shopDB.Close()

	log.Println("✅ Database connections successful.")

	// 4. Chạy seeder
	s := seeder.NewSeeder(productDB, shopDB)
	s.SeedProducts(*productCount)

	log.Println("🎉 Seeding complete.")
}
