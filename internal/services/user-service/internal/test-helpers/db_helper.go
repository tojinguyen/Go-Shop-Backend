// internal/services/user-service/internal/test_helpers/db_helper.go
package test_helpers

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
)

// SetupTestDatabase spins up a PostgreSQL container and runs migrations.
// It returns a connected DatabaseService and a cleanup function.
func SetupTestDatabase(t *testing.T) (*postgresql_infra.PostgreSQLService, func()) {
	ctx := context.Background()

	// 1. Create a PostgreSQL container
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Minute),
		),
	)
	if err != nil {
		t.Fatalf("could not start postgres container: %s", err)
	}

	// 2. Create the connection string for the container
	host, err := pgContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %s", err)
	}
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get container port: %s", err)
	}
	connStr := fmt.Sprintf("host=%s port=%s user=postgres password=password dbname=test-db sslmode=disable", host, port.Port())

	// 3. Connect to the new database
	dbConfig := &postgresql_infra.DatabaseConfig{
		Host:         host,
		Port:         port.Port(),
		User:         "postgres",
		Password:     "password",
		Name:         "test-db",
		SSLMode:      "disable",
		MaxOpenConns: 10,
		MaxIdleConns: 5,
	}

	dbService, err := postgresql_infra.NewPostgreSQLService(dbConfig)
	if err != nil {
		t.Fatalf("failed to connect to test database: %s", err)
	}

	// 4. Run Migrations using Goose
	goose.SetBaseFS(nil)
	// IMPORTANT: Adjust the path to your migrations directory
	migrationsDir := "../db/migrations"
	absMigrationsDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		t.Fatalf("could not get absolute path for migrations: %s", err)
	}

	db, err := goose.OpenDBWithDriver("postgres", connStr)
	if err != nil {
		t.Fatalf("goose: failed to open DB: %v\n", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	if err := goose.Up(db, absMigrationsDir); err != nil {
		t.Fatalf("goose: up failed: %v", err)
	}
	log.Println("Migrations applied successfully.")

	// 5. Define a cleanup function to terminate the container
	cleanup := func() {
		log.Println("Terminating test database container...")
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("could not terminate postgres container: %s", err)
		}
	}

	return dbService, cleanup
}
