package postgresql_infra

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-username/go-shop/internal/services/user-service/internal/db/sqlc"
)

// DatabaseService defines the interface for database operations
type DatabaseService interface {
	// Connection management
	GetQueries() *sqlc.Queries
	GetPool() *pgxpool.Pool
	GetConnection(ctx context.Context) (*pgxpool.Conn, error)

	// Transaction management
	BeginTransaction(ctx context.Context) (pgx.Tx, error)
	WithTransaction(ctx context.Context, fn func(*sqlc.Queries) error) error
	ExecuteInTransaction(ctx context.Context, queries []string) error

	// Lifecycle management
	Close()

	// User operations
	CreateUserAccount(ctx context.Context, email, hashedPassword string) (*sqlc.CreateUserAccountRow, error)
}

// Ensure PostgreSQLService implements DatabaseService interface
var _ DatabaseService = (*PostgreSQLService)(nil)
