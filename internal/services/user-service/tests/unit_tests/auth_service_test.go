package unittests

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	jwtService "github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/jwt"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/repository"
	"github.com/toji-dev/go-shop/internal/services/user-service/mocks"
)

// TestAuthService wraps AuthService for testing with injectable dependencies
type TestAuthService struct {
	userAccountRepo repository.UserAccountRepository
	jwtService      jwtService.JwtService
}

func NewTestAuthService(userRepo repository.UserAccountRepository, jwt jwtService.JwtService) *TestAuthService {
	return &TestAuthService{
		userAccountRepo: userRepo,
		jwtService:      jwt,
	}
}

func (s *TestAuthService) Register(ctx *gin.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
	// Check if user already exists
	exists, err := s.userAccountRepo.CheckUserExistsByEmail(ctx, req.Email)
	if err != nil {
		return dto.RegisterResponse{}, errors.New("failed to check user existence: " + err.Error())
	}

	if exists {
		return dto.RegisterResponse{}, errors.New("user with email " + req.Email + " already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.RegisterResponse{}, errors.New("failed to hash password: " + err.Error())
	}

	// Create user account parameters
	params := sqlc.CreateUserAccountParams{
		Email:          req.Email,
		HashedPassword: string(hashedPassword),
	}

	// Create the user account
	userAccount, err := s.userAccountRepo.CreateUserAccount(ctx, params)
	if err != nil {
		return dto.RegisterResponse{}, errors.New("failed to create user account: " + err.Error())
	}

	// Return success response
	return dto.RegisterResponse{
		UserID: userAccount.Id,
	}, nil
}

func (s *TestAuthService) Login(ctx *gin.Context, req dto.LoginRequest) (dto.AuthResponse, error) {
	// Get user by email
	userAccount, err := s.userAccountRepo.GetUserAccountByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return dto.AuthResponse{}, errors.New("invalid email or password")
		}
		return dto.AuthResponse{}, errors.New("failed to get user account: " + err.Error())
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(userAccount.HashedPassword), []byte(req.Password))
	if err != nil {
		return dto.AuthResponse{}, errors.New("invalid email or password")
	}

	// Generate JWT tokens
	tokenInput := &jwtService.GenerateTokenInput{
		UserId: userAccount.Id,
		Email:  userAccount.Email,
		Role:   userAccount.Role,
	}

	accessToken, err := s.jwtService.GenerateAccessToken(tokenInput)
	if err != nil {
		return dto.AuthResponse{}, errors.New("failed to generate access token: " + err.Error())
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(tokenInput)
	if err != nil {
		return dto.AuthResponse{}, errors.New("failed to generate refresh token: " + err.Error())
	}

	// Update last login time (optional, don't fail login if this fails)
	_ = s.userAccountRepo.UpdateLastLoginAt(ctx, userAccount.Id)

	// Return auth response
	return dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
		ID:           userAccount.Id,
		Email:        userAccount.Email,
		Role:         userAccount.Role,
	}, nil
}

// MockJWTService is a mock implementation of JWT service
type MockJWTService struct {
	mock.Mock
}

func NewMockJWTService() *MockJWTService {
	return &MockJWTService{}
}

func (m *MockJWTService) GenerateAccessToken(input *jwtService.GenerateTokenInput) (string, error) {
	args := m.Called(input)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GenerateRefreshToken(input *jwtService.GenerateTokenInput) (string, error) {
	args := m.Called(input)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateAccessToken(ctx context.Context, token string) (*jwtService.CustomJwtClaims, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwtService.CustomJwtClaims), args.Error(1)
}

func (m *MockJWTService) ValidateRefreshToken(token string) (*jwtService.CustomJwtClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwtService.CustomJwtClaims), args.Error(1)
}

// Helper function to create test user account with hashed password
func createTestUserAccountWithPassword(email, password string) *domain.UserAccount {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return &domain.UserAccount{
		Id:             "user-123",
		Email:          email,
		HashedPassword: string(hashedPassword),
		LastLoginAt:    "2025-01-01T00:00:00Z",
		Role:           "user",
		CreatedAt:      "2025-01-01T00:00:00Z",
		UpdatedAt:      "2025-01-01T00:00:00Z",
	}
}

// Helper function to create gin context for testing
func createTestGinContext() *gin.Context {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)
	return ctx
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name             string
		request          dto.RegisterRequest
		setupMocks       func(*mocks.MockUserAccountRepository)
		wantError        bool
		expectedErrorMsg string
		expectedUserID   string
	}{
		{
			name: "Success - User registered successfully",
			request: dto.RegisterRequest{
				Email:           "newuser@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			setupMocks: func(mockRepo *mocks.MockUserAccountRepository) {
				// Mock: CheckUserExistsByEmail returns false (user doesn't exist)
				mockRepo.On(
					"CheckUserExistsByEmail",
					mock.Anything,
					"newuser@example.com",
				).Return(false, nil)

				// Mock: CreateUserAccount returns new user
				expectedUser := &domain.UserAccount{
					Id:    "new-user-123",
					Email: "newuser@example.com",
					Role:  "user",
				}
				mockRepo.On(
					"CreateUserAccount",
					mock.Anything,
					mock.MatchedBy(func(params sqlc.CreateUserAccountParams) bool {
						return params.Email == "newuser@example.com"
					}),
				).Return(expectedUser, nil)
			},
			wantError:      false,
			expectedUserID: "new-user-123",
		},
		{
			name: "Error - User already exists",
			request: dto.RegisterRequest{
				Email:           "existing@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			setupMocks: func(mockRepo *mocks.MockUserAccountRepository) {
				// Mock: CheckUserExistsByEmail returns true (user exists)
				mockRepo.On(
					"CheckUserExistsByEmail",
					mock.Anything,
					"existing@example.com",
				).Return(true, nil)
			},
			wantError:        true,
			expectedErrorMsg: "user with email existing@example.com already exists",
		},
		{
			name: "Error - Database error when checking user existence",
			request: dto.RegisterRequest{
				Email:           "test@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			setupMocks: func(mockRepo *mocks.MockUserAccountRepository) {
				// Mock: CheckUserExistsByEmail returns database error
				mockRepo.On(
					"CheckUserExistsByEmail",
					mock.Anything,
					"test@example.com",
				).Return(false, errors.New("database connection failed"))
			},
			wantError:        true,
			expectedErrorMsg: "failed to check user existence",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := mocks.NewMockUserAccountRepository()
			mockJWT := NewMockJWTService()
			tt.setupMocks(mockRepo)

			// Create test service
			authService := NewTestAuthService(mockRepo, mockJWT)
			ginCtx := createTestGinContext()

			// Act
			result, err := authService.Register(ginCtx, tt.request)

			// Assert
			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				assert.Empty(t, result.UserID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserID, result.UserID)
			}

			// Verify all mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name             string
		request          dto.LoginRequest
		setupMocks       func(*mocks.MockUserAccountRepository, *MockJWTService)
		wantError        bool
		expectedErrorMsg string
		expectedEmail    string
		expectedRole     string
	}{
		{
			name: "Success - User login successfully",
			request: dto.LoginRequest{
				Email:    "user@example.com",
				Password: "password123",
			},
			setupMocks: func(mockRepo *mocks.MockUserAccountRepository, mockJWT *MockJWTService) {
				// Mock: GetUserAccountByEmail returns user
				user := createTestUserAccountWithPassword("user@example.com", "password123")
				mockRepo.On(
					"GetUserAccountByEmail",
					mock.Anything,
					"user@example.com",
				).Return(user, nil)

				// Mock: UpdateLastLoginAt succeeds
				mockRepo.On(
					"UpdateLastLoginAt",
					mock.Anything,
					user.Id,
				).Return(nil)

				// Mock: JWT token generation
				tokenInput := &jwtService.GenerateTokenInput{
					UserId: user.Id,
					Email:  user.Email,
					Role:   user.Role,
				}
				mockJWT.On("GenerateAccessToken", tokenInput).Return("access_token_123", nil)
				mockJWT.On("GenerateRefreshToken", tokenInput).Return("refresh_token_123", nil)
			},
			wantError:     false,
			expectedEmail: "user@example.com",
			expectedRole:  "user",
		},
		{
			name: "Error - User not found",
			request: dto.LoginRequest{
				Email:    "notfound@example.com",
				Password: "password123",
			},
			setupMocks: func(mockRepo *mocks.MockUserAccountRepository, mockJWT *MockJWTService) {
				// Mock: GetUserAccountByEmail returns not found error
				mockRepo.On(
					"GetUserAccountByEmail",
					mock.Anything,
					"notfound@example.com",
				).Return((*domain.UserAccount)(nil), sql.ErrNoRows)
			},
			wantError:        true,
			expectedErrorMsg: "invalid email or password",
		},
		{
			name: "Error - Wrong password",
			request: dto.LoginRequest{
				Email:    "user@example.com",
				Password: "wrongpassword",
			},
			setupMocks: func(mockRepo *mocks.MockUserAccountRepository, mockJWT *MockJWTService) {
				// Mock: GetUserAccountByEmail returns user with different password
				user := createTestUserAccountWithPassword("user@example.com", "correctpassword")
				mockRepo.On(
					"GetUserAccountByEmail",
					mock.Anything,
					"user@example.com",
				).Return(user, nil)
			},
			wantError:        true,
			expectedErrorMsg: "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := mocks.NewMockUserAccountRepository()
			mockJWT := NewMockJWTService()
			tt.setupMocks(mockRepo, mockJWT)

			authService := NewTestAuthService(mockRepo, mockJWT)
			ginCtx := createTestGinContext()

			// Act
			result, err := authService.Login(ginCtx, tt.request)

			// Assert
			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				assert.Empty(t, result.AccessToken)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
				assert.Equal(t, "Bearer", result.TokenType)
				assert.Equal(t, int64(3600), result.ExpiresIn)
				assert.Equal(t, tt.expectedEmail, result.Email)
				assert.Equal(t, tt.expectedRole, result.Role)
			}

			// Verify all mock expectations
			mockRepo.AssertExpectations(t)
			if !tt.wantError || !assert.ObjectsAreEqual(tt.expectedErrorMsg, "invalid email or password") {
				mockJWT.AssertExpectations(t)
			}
		})
	}
}
