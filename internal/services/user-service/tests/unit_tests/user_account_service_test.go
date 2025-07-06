package unittests

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/user-service/mocks"
)

func createTestUserAccount() *domain.UserAccount {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	return &domain.UserAccount{
		Id:             "user-123",
		Email:          "test@example.com",
		HashedPassword: string(hashedPassword),
		LastLoginAt:    "2025-01-01T00:00:00Z",
		Role:           "user",
		CreatedAt:      "2025-01-01T00:00:00Z",
		UpdatedAt:      "2025-01-01T00:00:00Z",
	}
}

func TestUserAccountRepository_CreateUserAccount(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		params        sqlc.CreateUserAccountParams
		mockReturn    *domain.UserAccount
		mockError     error
		wantError     bool
		expectedEmail string
	}{
		{
			name: "Success - Create user account",
			params: sqlc.CreateUserAccountParams{
				Email:          "test@example.com",
				HashedPassword: "hashed_password",
			},
			mockReturn:    createTestUserAccount(),
			mockError:     nil,
			wantError:     false,
			expectedEmail: "test@example.com",
		},
		{
			name: "Error - Repository failure",
			params: sqlc.CreateUserAccountParams{
				Email:          "test@example.com",
				HashedPassword: "hashed_password",
			},
			mockReturn: nil,
			mockError:  errors.New("database connection failed"),
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := mocks.NewMockUserAccountRepository()
			mockRepo.On("CreateUserAccount", ctx, tt.params).Return(tt.mockReturn, tt.mockError)

			// Act
			result, err := mockRepo.CreateUserAccount(ctx, tt.params)

			// Assert
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tt.mockError, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedEmail, result.Email)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserAccountRepository_GetUserAccountByEmail(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		email        string
		mockReturn   *domain.UserAccount
		mockError    error
		wantError    bool
		expectedUser *domain.UserAccount
	}{
		{
			name:         "Success - User found",
			email:        "test@example.com",
			mockReturn:   createTestUserAccount(),
			mockError:    nil,
			wantError:    false,
			expectedUser: createTestUserAccount(),
		},
		{
			name:       "Error - User not found",
			email:      "notfound@example.com",
			mockReturn: nil,
			mockError:  sql.ErrNoRows,
			wantError:  true,
		},
		{
			name:       "Error - Database error",
			email:      "test@example.com",
			mockReturn: nil,
			mockError:  errors.New("database connection lost"),
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := mocks.NewMockUserAccountRepository()
			mockRepo.On("GetUserAccountByEmail", ctx, tt.email).Return(tt.mockReturn, tt.mockError)

			// Act
			result, err := mockRepo.GetUserAccountByEmail(ctx, tt.email)

			// Assert
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tt.mockError, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedUser.Email, result.Email)
				assert.Equal(t, tt.expectedUser.Id, result.Id)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserAccountRepository_GetUserAccountByID(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - User found by ID", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewMockUserAccountRepository()
		userID := "user-123"
		expectedUser := createTestUserAccount()
		mockRepo.On("GetUserAccountByID", ctx, userID).Return(expectedUser, nil)

		// Act
		result, err := mockRepo.GetUserAccountByID(ctx, userID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedUser.Id, result.Id)
		assert.Equal(t, expectedUser.Email, result.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - User not found", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewMockUserAccountRepository()
		userID := "non-existent-id"
		mockRepo.On("GetUserAccountByID", ctx, userID).Return((*domain.UserAccount)(nil), sql.ErrNoRows)

		// Act
		result, err := mockRepo.GetUserAccountByID(ctx, userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, sql.ErrNoRows, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserAccountRepository_UpdateLastLoginAt(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - Update last login", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewMockUserAccountRepository()
		userID := "user-123"
		mockRepo.On("UpdateLastLoginAt", ctx, userID).Return(nil)

		// Act
		err := mockRepo.UpdateLastLoginAt(ctx, userID)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Update failed", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewMockUserAccountRepository()
		userID := "user-123"
		expectedError := errors.New("update failed")
		mockRepo.On("UpdateLastLoginAt", ctx, userID).Return(expectedError)

		// Act
		err := mockRepo.UpdateLastLoginAt(ctx, userID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserAccountRepository_UpdateUserPassword(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - Update password", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewMockUserAccountRepository()
		userID := "user-123"
		hashedPassword := "new_hashed_password"
		mockRepo.On("UpdateUserPassword", ctx, userID, hashedPassword).Return(nil)

		// Act
		err := mockRepo.UpdateUserPassword(ctx, userID, hashedPassword)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Update failed", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewMockUserAccountRepository()
		userID := "user-123"
		hashedPassword := "new_hashed_password"
		expectedError := errors.New("password update failed")
		mockRepo.On("UpdateUserPassword", ctx, userID, hashedPassword).Return(expectedError)

		// Act
		err := mockRepo.UpdateUserPassword(ctx, userID, hashedPassword)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserAccountRepository_SoftDeleteUserAccount(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - Soft delete user", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewMockUserAccountRepository()
		userID := "user-123"
		mockRepo.On("SoftDeleteUserAccount", ctx, userID).Return(nil)

		// Act
		err := mockRepo.SoftDeleteUserAccount(ctx, userID)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Delete failed", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewMockUserAccountRepository()
		userID := "user-123"
		expectedError := errors.New("soft delete failed")
		mockRepo.On("SoftDeleteUserAccount", ctx, userID).Return(expectedError)

		// Act
		err := mockRepo.SoftDeleteUserAccount(ctx, userID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserAccountRepository_CheckUserExistsByEmail(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		email         string
		mockReturn    bool
		mockError     error
		wantError     bool
		expectedExist bool
	}{
		{
			name:          "Success - User exists",
			email:         "existing@example.com",
			mockReturn:    true,
			mockError:     nil,
			wantError:     false,
			expectedExist: true,
		},
		{
			name:          "Success - User does not exist",
			email:         "nonexistent@example.com",
			mockReturn:    false,
			mockError:     nil,
			wantError:     false,
			expectedExist: false,
		},
		{
			name:       "Error - Database error",
			email:      "test@example.com",
			mockReturn: false,
			mockError:  errors.New("database query failed"),
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := mocks.NewMockUserAccountRepository()
			mockRepo.On("CheckUserExistsByEmail", ctx, tt.email).Return(tt.mockReturn, tt.mockError)

			// Act
			exists, err := mockRepo.CheckUserExistsByEmail(ctx, tt.email)

			// Assert
			if tt.wantError {
				assert.Error(t, err)
				assert.Equal(t, tt.mockError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedExist, exists)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

// Integration test example showing how the mock would be used in service testing
func TestUserAccountRepository_IntegrationExample(t *testing.T) {
	ctx := context.Background()

	t.Run("Complete user registration flow", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewMockUserAccountRepository()

		// Step 1: Check if user exists (should return false)
		email := "newuser@example.com"
		mockRepo.On("CheckUserExistsByEmail", ctx, email).Return(false, nil)

		// Step 2: Create user account
		params := sqlc.CreateUserAccountParams{
			Email:          email,
			HashedPassword: "hashed_password",
		}
		newUser := &domain.UserAccount{
			Id:             "new-user-123",
			Email:          email,
			HashedPassword: "hashed_password",
			Role:           "user",
			CreatedAt:      "2025-01-01T00:00:00Z",
			UpdatedAt:      "2025-01-01T00:00:00Z",
		}
		mockRepo.On("CreateUserAccount", ctx, params).Return(newUser, nil)

		// Act - Simulate registration flow
		exists, err := mockRepo.CheckUserExistsByEmail(ctx, email)
		assert.NoError(t, err)
		assert.False(t, exists)

		createdUser, err := mockRepo.CreateUserAccount(ctx, params)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, email, createdUser.Email)

		// Assert all expectations met
		mockRepo.AssertExpectations(t)
	})
}
