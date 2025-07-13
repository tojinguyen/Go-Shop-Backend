package services_test

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	email_mocks "github.com/toji-dev/go-shop/internal/pkg/email/mocks"
	redis_mocks "github.com/toji-dev/go-shop/internal/pkg/infra/redis-infra/mocks"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	jwt_mocks "github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/jwt/mocks"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/repository"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/services"
	test_helpers "github.com/toji-dev/go-shop/internal/services/user-service/internal/test-helpers"
)

func TestAuthService_Register_Integration(t *testing.T) {
	// Step 1: Setup the test database using our helper
	// This gives us a clean, migrated database for our test.
	dbService, cleanup := test_helpers.SetupTestDatabase(t)
	// The cleanup function will be called automatically at the end of the test.
	t.Cleanup(cleanup)

	// Step 2: Setup dependencies for the AuthService
	// We use the real repository connected to our test database.
	// Other dependencies like JWT, Redis, Email can be mocked for this specific test.
	userRepo := repository.NewUserAccountRepository(dbService)
	// Mocks for services we are not testing directly
	// (for a pure Register test, we don't need them, but it's good practice to have them)
	mockJWT := &jwt_mocks.JwtService{}
	mockRedis := &redis_mocks.RedisServiceInterface{}
	mockEmail := &email_mocks.EmailService{}
	cfg := &config.Config{} // Empty config is fine for this test

	authService := services.NewAuthService(userRepo, mockJWT, mockRedis, mockEmail, cfg)

	// Step 3: Define test cases
	t.Run("Success - Register a new user", func(t *testing.T) {
		// Arrange
		ctx, _ := gin.CreateTestContext(nil)
		req := dto.RegisterRequest{
			Email:           "test.user@example.com",
			Password:        "Password123!",
			ConfirmPassword: "Password123!",
		}

		// Act
		resp, err := authService.Register(ctx, req)

		// Assert
		assert.NoError(t, err)
		require.NotNil(t, resp)
		assert.NotEmpty(t, resp.UserID)

		// **Verification Step**: Directly query the database to confirm the user was created.
		// This is the core of an integration test.
		var dbEmail, dbPassword string
		query := "SELECT email, hashed_password FROM user_accounts WHERE id = $1"
		err = dbService.GetPool().QueryRow(context.Background(), query, resp.UserID).Scan(&dbEmail, &dbPassword)

		require.NoError(t, err, "User should exist in the database after registration")
		assert.Equal(t, "test.user@example.com", dbEmail)
		assert.NotEmpty(t, dbPassword, "Password should be hashed and stored")
	})

	t.Run("Failure - Attempt to register an existing user", func(t *testing.T) {
		// Arrange: First, create a user directly in the DB to simulate an existing user.
		existingEmail := "existing.user@example.com"
		insertUser(t, dbService.GetPool(), existingEmail)

		ctx, _ := gin.CreateTestContext(nil)
		req := dto.RegisterRequest{
			Email:    existingEmail,
			Password: "Password123!",
		}

		// Act
		resp, err := authService.Register(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "user with email existing.user@example.com already exists")
	})
}

// Helper function to pre-seed data for tests
func insertUser(t *testing.T, pool *pgxpool.Pool, email string) {
	t.Helper()
	query := "INSERT INTO user_accounts (email, hashed_password, user_role) VALUES ($1, $2, 'customer')"
	_, err := pool.Exec(context.Background(), query, email, "somehashedpassword")
	require.NoError(t, err)
}
