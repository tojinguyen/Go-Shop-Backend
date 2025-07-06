package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/repository"
)

// MockUserAccountRepository là mock implementation của UserAccountRepository interface
// Sử dụng testify/mock framework để tạo mock object
type MockUserAccountRepository struct {
	mock.Mock
}

// NewMockUserAccountRepository tạo một instance mới của mock repository
func NewMockUserAccountRepository() *MockUserAccountRepository {
	return &MockUserAccountRepository{}
}

// CreateUserAccount mock implementation với testify/mock
func (m *MockUserAccountRepository) CreateUserAccount(ctx context.Context, params sqlc.CreateUserAccountParams) (*domain.UserAccount, error) {
	args := m.Called(ctx, params)

	// Kiểm tra nếu return value đầu tiên là nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.UserAccount), args.Error(1)
}

// GetUserAccountByEmail mock implementation với testify/mock
func (m *MockUserAccountRepository) GetUserAccountByEmail(ctx context.Context, email string) (*domain.UserAccount, error) {
	args := m.Called(ctx, email)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.UserAccount), args.Error(1)
}

// GetUserAccountByID mock implementation với testify/mock
func (m *MockUserAccountRepository) GetUserAccountByID(ctx context.Context, id string) (*domain.UserAccount, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.UserAccount), args.Error(1)
}

// UpdateLastLoginAt mock implementation với testify/mock
func (m *MockUserAccountRepository) UpdateLastLoginAt(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// UpdateUserPassword mock implementation với testify/mock
func (m *MockUserAccountRepository) UpdateUserPassword(ctx context.Context, id string, hashedPassword string) error {
	args := m.Called(ctx, id, hashedPassword)
	return args.Error(0)
}

// SoftDeleteUserAccount mock implementation với testify/mock
func (m *MockUserAccountRepository) SoftDeleteUserAccount(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// CheckUserExistsByEmail mock implementation với testify/mock
func (m *MockUserAccountRepository) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

// Verify that MockUserAccountRepository implements repository.UserAccountRepository
var _ repository.UserAccountRepository = (*MockUserAccountRepository)(nil)
