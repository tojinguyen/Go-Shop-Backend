package main

import (
	"flag"
	"log"

	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/seeder"
)

func main() {
	// Define a command-line flag for the number of users to create
	count := flag.Int("count", 10, "Number of users to seed")
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

	// 3. Run the seeder function
	seeder.SeedUsers(dbService.GetPool(), *count)

	log.Println("Seeding complete.")
}
