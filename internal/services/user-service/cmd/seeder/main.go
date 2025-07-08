// internal/services/user-service/cmd/seeder/main.go
package main

import (
	"flag"
	"log"

	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/seeder"
)

func main() {
	// Define command-line flags
	userCount := flag.Int("users", 10, "Number of regular customer users to seed")
	shipperCount := flag.Int("shippers", 5, "Number of shipper users to seed")
	flag.Parse()

	log.Println("Starting database seeder...")

	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration loaded successfully.")

	// 2. Connect to the database
	dbConfig := &postgresql_infra.DatabaseConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
	}

	dbService, err := postgresql_infra.NewPostgreSQLService(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbService.Close()
	log.Println("Database connection successful.")

	// 3. Run the seeder
	s := seeder.NewSeeder(dbService.GetPool())
	s.SeedAll(*userCount, *shipperCount)

	log.Println("Seeding complete.")
}
