package repository

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/constant"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/domain"
)

type UserAccountRepository interface {
	CreateUserAccount(ctx context.Context, params sqlc.CreateUserAccountParams) (*domain.UserAccount, error)
	GetUserAccountByEmail(ctx context.Context, email string) (*domain.UserAccount, error)
	GetUserAccountByID(ctx context.Context, id string) (*domain.UserAccount, error)
	UpdateLastLoginAt(ctx context.Context, id string) error
	UpdateUserPassword(ctx context.Context, id string, hashedPassword string) error
	SoftDeleteUserAccount(ctx context.Context, id string) error
	CheckUserExistsByEmail(ctx context.Context, email string) (bool, error)
}

type userAccountRepository struct {
	db      *postgresql_infra.PostgreSQLService
	queries *sqlc.Queries
}

// NewUserAccountRepository creates a new instance of UserAccountRepository
func NewUserAccountRepository(db *postgresql_infra.PostgreSQLService) UserAccountRepository {
	return &userAccountRepository{
		db:      db,
		queries: sqlc.New(db.GetPool()),
	}
}

// CreateUserAccount creates a new user account
func (r *userAccountRepository) CreateUserAccount(ctx context.Context, params sqlc.CreateUserAccountParams) (*domain.UserAccount, error) {
	// Use SQLC generated function
	sqlcParams := sqlc.CreateUserAccountParams{
		Email:          params.Email,
		HashedPassword: params.HashedPassword,
		UserRole:       string(constant.UserRoleCustomer),
	}

	result, err := r.queries.CreateUserAccount(ctx, sqlcParams)
	if err != nil {
		log.Println("Error creating user account:", err)
		return nil, err
	}

	// Convert SQLC result to domain model
	userAccount := &domain.UserAccount{
		Id:             result.ID.String(),
		Email:          result.Email,
		HashedPassword: "", // Don't return password in create response
		LastLoginAt:    "", // Will be empty for new accounts
		CreatedAt:      result.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:      result.UpdatedAt.Time.Format(time.RFC3339),
	}

	return userAccount, nil
}

// GetUserAccountByEmail retrieves a user account by email
func (r *userAccountRepository) GetUserAccountByEmail(ctx context.Context, email string) (*domain.UserAccount, error) {
	// Use SQLC generated function
	result, err := r.queries.GetUserAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Convert SQLC result to domain model
	userAccount := &domain.UserAccount{
		Id:             result.ID.String(),
		Email:          result.Email,
		HashedPassword: result.HashedPassword,
		Role:           result.UserRole,
	}

	// Handle nullable timestamp fields
	if result.LastLoginAt.Valid {
		userAccount.LastLoginAt = result.LastLoginAt.Time.Format(time.RFC3339)
	}
	userAccount.CreatedAt = result.CreatedAt.Time.Format(time.RFC3339)
	userAccount.UpdatedAt = result.UpdatedAt.Time.Format(time.RFC3339)

	return userAccount, nil
}

// GetUserAccountByID retrieves a user account by ID
func (r *userAccountRepository) GetUserAccountByID(ctx context.Context, id string) (*domain.UserAccount, error) {
	// Convert string ID to UUID
	uuid := pgtype.UUID{}
	err := uuid.Scan(id)
	if err != nil {
		return nil, err
	}

	// Use SQLC generated function
	result, err := r.queries.GetUserAccountByID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	// Convert SQLC result to domain model
	userAccount := &domain.UserAccount{
		Id:    result.ID.String(),
		Email: result.Email,
	}

	// Handle nullable timestamp fields
	if result.LastLoginAt.Valid {
		userAccount.LastLoginAt = result.LastLoginAt.Time.Format(time.RFC3339)
	}
	userAccount.CreatedAt = result.CreatedAt.Time.Format(time.RFC3339)
	userAccount.UpdatedAt = result.UpdatedAt.Time.Format(time.RFC3339)

	return userAccount, nil
}

// UpdateLastLoginAt updates the last login timestamp for a user
func (r *userAccountRepository) UpdateLastLoginAt(ctx context.Context, id string) error {
	// Convert string ID to UUID
	uuid := pgtype.UUID{}
	err := uuid.Scan(id)
	if err != nil {
		return err
	}

	// Use SQLC generated function
	return r.queries.UpdateLastLoginAt(ctx, uuid)
}

// UpdateUserPassword updates the password for a user account
func (r *userAccountRepository) UpdateUserPassword(ctx context.Context, id string, hashedPassword string) error {
	// Convert string ID to UUID
	uuid := pgtype.UUID{}
	err := uuid.Scan(id)
	if err != nil {
		return err
	}

	// Use SQLC generated function
	params := sqlc.UpdateUserPasswordParams{
		ID:             uuid,
		HashedPassword: hashedPassword,
	}

	return r.queries.UpdateUserPassword(ctx, params)
}

// SoftDeleteUserAccount soft deletes a user account
func (r *userAccountRepository) SoftDeleteUserAccount(ctx context.Context, id string) error {
	// Convert string ID to UUID
	uuid := pgtype.UUID{}
	err := uuid.Scan(id)
	if err != nil {
		return err
	}

	// Use SQLC generated function
	return r.queries.SoftDeleteUserAccount(ctx, uuid)
}

// Check if user exists by email
func (r *userAccountRepository) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	// Use SQLC generated function
	id, err := r.queries.CheckUserExistsByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		log.Println("Error checking user existence by email:", err)
		return false, err
	}

	// If id is not null, user exists
	return id != pgtype.UUID{}, nil
}
