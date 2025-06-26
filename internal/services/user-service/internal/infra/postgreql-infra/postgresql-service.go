package postgresql_infra

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-username/go-shop/internal/services/user-service/internal/config"
	"github.com/your-username/go-shop/internal/services/user-service/internal/db/sqlc"
)

// PostgreSQLService represents the PostgreSQL database service
type PostgreSQLService struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
	config  *config.DatabaseConfig
}

// NewPostgreSQLService creates a new PostgreSQL service instance
func NewPostgreSQLService(cfg *config.DatabaseConfig) (*PostgreSQLService, error) {
	service := &PostgreSQLService{
		config: cfg,
	}

	if err := service.connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize SQLC queries
	service.queries = sqlc.New(service.pool)

	return service, nil
}

// connect establishes connection to PostgreSQL database
func (s *PostgreSQLService) connect() error {
	// Build connection string
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		s.config.User,
		s.config.Password,
		s.config.Host,
		s.config.Port,
		s.config.Name,
		s.config.SSLMode,
	)

	// Configure connection pool
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse database config: %w", err)
	}

	// Set connection pool settings
	poolConfig.MaxConns = int32(s.config.MaxOpenConns)
	poolConfig.MinConns = int32(s.config.MaxIdleConns)
	poolConfig.MaxConnLifetime = s.config.MaxLifetime
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	// Create connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	s.pool = pool
	log.Printf("Successfully connected to PostgreSQL database: %s:%s/%s",
		s.config.Host, s.config.Port, s.config.Name)

	return nil
}

// GetQueries returns the SQLC queries instance
func (s *PostgreSQLService) GetQueries() *sqlc.Queries {
	return s.queries
}

// GetPool returns the connection pool
func (s *PostgreSQLService) GetPool() *pgxpool.Pool {
	return s.pool
}

// GetConnection gets a single connection from the pool
func (s *PostgreSQLService) GetConnection(ctx context.Context) (*pgxpool.Conn, error) {
	return s.pool.Acquire(ctx)
}

// BeginTransaction starts a new database transaction
func (s *PostgreSQLService) BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	return s.pool.Begin(ctx)
}

// WithTransaction executes a function within a database transaction
func (s *PostgreSQLService) WithTransaction(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := s.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Create queries with transaction
	qtx := s.queries.WithTx(tx)

	// Execute function
	if err := fn(qtx); err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Close closes the database connection pool
func (s *PostgreSQLService) Close() {
	if s.pool != nil {
		log.Println("Closing PostgreSQL connection pool...")
		s.pool.Close()
	}
}

// ExecuteInTransaction executes multiple queries in a single transaction
func (s *PostgreSQLService) ExecuteInTransaction(ctx context.Context, queries []string) error {
	tx, err := s.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for i, query := range queries {
		if _, err := tx.Exec(ctx, query); err != nil {
			return fmt.Errorf("failed to execute query %d: %w", i+1, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
