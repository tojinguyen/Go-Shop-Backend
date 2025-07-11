package services_test

import (
	"database/sql"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	// Import các mock đã được mockery sinh ra
	email_mocks "github.com/toji-dev/go-shop/internal/pkg/email/mocks"
	redis_mocks "github.com/toji-dev/go-shop/internal/pkg/infra/redis-infra/mocks"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	jwt_pkg "github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/jwt"
	jwt_mocks "github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/jwt/mocks"
	repo_mocks "github.com/toji-dev/go-shop/internal/services/user-service/internal/repository/mocks"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/services"
)

// Helper function không đổi
func createTestGinContext() *gin.Context {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)
	return ctx
}

func TestAuthService_Register(t *testing.T) {
	testCases := []struct {
		name          string
		request       dto.RegisterRequest
		setupMocks    func(mockRepo *repo_mocks.UserAccountRepository)
		expectedID    string
		expectError   bool
		expectedError string
	}{
		{
			name:    "Success - User registered successfully",
			request: dto.RegisterRequest{Email: "newuser@example.com", Password: "Password123!"},
			setupMocks: func(mockRepo *repo_mocks.UserAccountRepository) {
				// Sử dụng EXPECT() thay vì On()
				mockRepo.EXPECT().CheckUserExistsByEmail(mock.Anything, "newuser@example.com").Return(false, nil).Once()

				createdUser := &domain.UserAccount{Id: "new-user-id", Email: "newuser@example.com"}
				mockRepo.EXPECT().CreateUserAccount(mock.Anything, mock.AnythingOfType("sqlc.CreateUserAccountParams")).Return(createdUser, nil).Once()
			},
			expectedID:  "new-user-id",
			expectError: false,
		},
		{
			name:    "Error - User already exists",
			request: dto.RegisterRequest{Email: "existing@example.com", Password: "Password123!"},
			setupMocks: func(mockRepo *repo_mocks.UserAccountRepository) {
				mockRepo.EXPECT().CheckUserExistsByEmail(mock.Anything, "existing@example.com").Return(true, nil).Once()
			},
			expectError:   true,
			expectedError: "user with email existing@example.com already exists",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo := new(repo_mocks.UserAccountRepository)
			tc.setupMocks(mockUserRepo)

			// Khởi tạo service với các dependency đã được mock
			authService := services.NewAuthService(mockUserRepo, nil, nil, nil, &config.Config{})

			response, err := authService.Register(createTestGinContext(), tc.request)

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedID, response.UserID)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

// =================================================================
// Test hàm Login với cú pháp EXPECT()
// =================================================================
func TestAuthService_Login_WithExpect(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctPassword123!"), bcrypt.DefaultCost)

	testCases := []struct {
		name          string
		request       dto.LoginRequest
		setupMocks    func(mockRepo *repo_mocks.UserAccountRepository, mockJWT *jwt_mocks.JwtService)
		expectedResp  *dto.AuthResponse
		expectError   bool
		expectedError string
	}{
		{
			name:    "Success - Valid login credentials",
			request: dto.LoginRequest{Email: "user@example.com", Password: "correctPassword123!"},
			setupMocks: func(mockRepo *repo_mocks.UserAccountRepository, mockJWT *jwt_mocks.JwtService) {
				user := &domain.UserAccount{
					Id:             "user-123",
					Email:          "user@example.com",
					HashedPassword: string(hashedPassword),
					Role:           "customer",
				}
				mockRepo.EXPECT().GetUserAccountByEmail(mock.Anything, "user@example.com").Return(user, nil).Once()
				mockRepo.EXPECT().UpdateLastLoginAt(mock.Anything, "user-123").Return(nil).Once()

				tokenInput := &jwt_pkg.GenerateTokenInput{UserId: user.Id, Email: user.Email, Role: user.Role}
				mockJWT.EXPECT().GenerateAccessToken(tokenInput).Return("fake-access-token", nil).Once()
				mockJWT.EXPECT().GenerateRefreshToken(tokenInput).Return("fake-refresh-token", nil).Once()
			},
			expectedResp: &dto.AuthResponse{
				AccessToken:  "fake-access-token",
				RefreshToken: "fake-refresh-token",
				TokenType:    "Bearer",
				ExpiresIn:    3600,
				ID:           "user-123",
				Email:        "user@example.com",
				Role:         "customer",
			},
			expectError: false,
		},
		{
			name:    "Error - User not found",
			request: dto.LoginRequest{Email: "notfound@example.com", Password: "password"},
			setupMocks: func(mockRepo *repo_mocks.UserAccountRepository, mockJWT *jwt_mocks.JwtService) {
				mockRepo.EXPECT().GetUserAccountByEmail(mock.Anything, "notfound@example.com").Return(nil, sql.ErrNoRows).Once()
			},
			expectError:   true,
			expectedError: "invalid email or password",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// --- ARRANGE ---
			mockUserRepo := new(repo_mocks.UserAccountRepository)
			mockJWT := new(jwt_mocks.JwtService)
			tc.setupMocks(mockUserRepo, mockJWT)

			authService := services.NewAuthService(
				mockUserRepo, mockJWT, new(redis_mocks.RedisServiceInterface), new(email_mocks.EmailService), &config.Config{},
			)

			// --- ACT ---
			response, err := authService.Login(createTestGinContext(), tc.request)

			// --- ASSERT ---
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, *tc.expectedResp, response)
			}

			mockUserRepo.AssertExpectations(t)
			mockJWT.AssertExpectations(t)
		})
	}
}
