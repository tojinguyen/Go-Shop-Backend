package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	jwtService "github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/jwt"
)

type AuthService struct {
	container *container.ServiceContainer
}

func NewAuthService(container *container.ServiceContainer) *AuthService {
	return &AuthService{
		container: container,
	}
}

func (s *AuthService) Register(ctx *gin.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
	// Check if user already exists
	exists, err := s.container.GetUserAccountRepo().CheckUserExistsByEmail(context.Background(), req.Email)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("failed to check user existence: %w", err)
	}

	if exists {
		return dto.RegisterResponse{}, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user account parameters
	params := sqlc.CreateUserAccountParams{
		Email:          req.Email,
		HashedPassword: string(hashedPassword),
	}

	// Create the user account
	userAccount, err := s.container.GetUserAccountRepo().CreateUserAccount(context.Background(), params)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("failed to create user account: %w", err)
	}

	// Return success response
	return dto.RegisterResponse{
		UserID: userAccount.Id,
	}, nil
}

// Register function signature for backward compatibility
func Register(ctx *gin.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
	// This function would need a service container instance
	// It's better to use the AuthService struct method above
	return dto.RegisterResponse{}, fmt.Errorf("use AuthService.Register method instead")
}

func (s *AuthService) Login(ctx *gin.Context, req dto.LoginRequest) (dto.AuthResponse, error) {
	// Get user by email
	userAccount, err := s.container.GetUserAccountRepo().GetUserAccountByEmail(context.Background(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return dto.AuthResponse{}, fmt.Errorf("invalid email or password")
		}
		return dto.AuthResponse{}, fmt.Errorf("failed to get user account: %w", err)
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(userAccount.HashedPassword), []byte(req.Password))
	if err != nil {
		return dto.AuthResponse{}, fmt.Errorf("invalid email or password")
	}

	// Generate JWT tokens
	tokenInput := &jwtService.GenerateTokenInput{
		UserId: userAccount.Id,
		Email:  userAccount.Email,
		Role:   "user", // Default role
	}

	accessToken, err := s.container.GetJWT().GenerateAccessToken(tokenInput)
	if err != nil {
		return dto.AuthResponse{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.container.GetJWT().GenerateRefreshToken(tokenInput)
	if err != nil {
		return dto.AuthResponse{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Update last login time
	err = s.container.GetUserAccountRepo().UpdateLastLoginAt(context.Background(), userAccount.Id)
	if err != nil {
		// Log error but don't fail the login
		fmt.Printf("Failed to update last login time: %v\n", err)
	}

	// Return auth response
	return dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
		User: &dto.UserInfo{
			ID:       userAccount.Id,
			Email:    userAccount.Email,
			Username: userAccount.Email, // Using email as username for now
			Role:     "user",            // Default role
		},
	}, nil
}

func (s *AuthService) RefreshToken(ctx *gin.Context, req dto.RefreshTokenRequest) (dto.AuthResponse, error) {
	// Validate refresh token
	claims, err := s.container.GetJWT().ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return dto.AuthResponse{}, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get user account to ensure user still exists
	userAccount, err := s.container.GetUserAccountRepo().GetUserAccountByID(context.Background(), claims.UserId)
	if err != nil {
		return dto.AuthResponse{}, fmt.Errorf("user not found: %w", err)
	}

	// Generate new tokens
	tokenInput := &jwtService.GenerateTokenInput{
		UserId: userAccount.Id,
		Email:  userAccount.Email,
		Role:   claims.Role,
	}

	accessToken, err := s.container.GetJWT().GenerateAccessToken(tokenInput)
	if err != nil {
		return dto.AuthResponse{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.container.GetJWT().GenerateRefreshToken(tokenInput)
	if err != nil {
		return dto.AuthResponse{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Return new auth response
	return dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
		User: &dto.UserInfo{
			ID:       userAccount.Id,
			Email:    userAccount.Email,
			Username: userAccount.Email,
			Role:     claims.Role,
		},
	}, nil
}

func (s *AuthService) ValidateToken(ctx *gin.Context, token string) (*dto.TokenValidationResponse, error) {
	// Validate access token
	claims, err := s.container.GetJWT().ValidateAccessToken(ctx.Request.Context(), token)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	// Return token information
	tokenInfo := &dto.TokenValidationResponse{
		Valid:     true,
		UserID:    claims.UserId,
		UserEmail: claims.Email,
		UserRole:  claims.Role,
		ExpiresAt: claims.ExpiresAt,
	}

	return tokenInfo, nil
}
