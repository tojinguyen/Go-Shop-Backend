package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/db/seeder"
)

func main() {
	shopCount := flag.Int("shops", 50, "Number of shops to seed")
	flag.Parse()

	log.Println("🌱 Starting shop-service database seeder...")

	// 1. Tải cấu hình của shop-service
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// 2. Kết nối đến DB của shop-service
	shopDB, err := pgxpool.New(context.Background(), cfg.GetDatabaseURL())
	if err != nil {
		log.Fatalf("❌ Failed to connect to shop-service database: %v", err)
	}
	defer shopDB.Close()

	// 3. Kết nối đến DB của user-service (cần thông tin từ env)
	userDBConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.UserServiceDB.User,
		cfg.UserServiceDB.Password,
		cfg.UserServiceDB.Host,
		cfg.UserServiceDB.Port,
		cfg.UserServiceDB.DBName,
		cfg.UserServiceDB.SSLMode,
	)
	userDB, err := pgxpool.New(context.Background(), userDBConnStr)
	if err != nil {
		log.Fatalf("❌ Failed to connect to user-service database: %v", err)
	}
	defer userDB.Close()

	log.Println("✅ Database connections successful.")

	// 4. Chạy seeder
	s := seeder.NewSeeder(shopDB, userDB)
	s.SeedShops(*shopCount)

	log.Println("🎉 Seeding complete.")
}
