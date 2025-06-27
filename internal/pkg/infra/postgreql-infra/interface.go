package postgresql_infra

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseService defines the interface for database operations
type DatabaseService interface {
	// Connection management
	GetPool() *pgxpool.Pool
	GetConnection(ctx context.Context) (*pgxpool.Conn, error)

	// Transaction management
	BeginTransaction(ctx context.Context) (pgx.Tx, error)
	WithTransaction(ctx context.Context, fn func(pgx.Tx) error) error
	ExecuteInTransaction(ctx context.Context, queries []string) error

	// Lifecycle management
	Close()
}

// Ensure PostgreSQLService implements DatabaseService interface
var _ DatabaseService = (*PostgreSQLService)(nil)
