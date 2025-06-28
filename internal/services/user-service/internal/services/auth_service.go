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
	userAccount, err := s.container.GetUserAccountRepo().CreateUserAccount(ctx, params)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("failed to create user account: %w", err)
	}

	// Return success response
	return dto.RegisterResponse{
		UserID: userAccount.Id,
	}, nil
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
		Role:   userAccount.Role,
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
	err = s.container.GetUserAccountRepo().UpdateLastLoginAt(ctx, userAccount.Id)
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
		ID:           userAccount.Id,
		Email:        userAccount.Email,
		Role:         userAccount.Role,
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
		ID:           userAccount.Id,
		Email:        userAccount.Email,
		Role:         userAccount.Role,
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

func (s *AuthService) Logout(ctx *gin.Context, token string) error {
	// Extract user ID from token for logging purposes
	claims, err := s.container.GetJWT().ValidateAccessToken(ctx.Request.Context(), token)
	if err != nil {
		// Even if token is invalid, we consider logout successful
		return nil
	}

	// TODO: Implement token blacklisting in Redis
	// For now, we'll just return success since JWT tokens are stateless
	// In a production environment, you would:
	// 1. Add token to Redis blacklist with expiry time
	// 2. Log the logout event
	// 3. Clear any user sessions

	// Log the logout
	fmt.Printf("User %s logged out successfully\n", claims.UserId)

	return nil
}

func (s *AuthService) ForgotPassword(ctx *gin.Context, req dto.ForgotPasswordRequest) error {
	// Check if user exists
	userAccount, err := s.container.GetUserAccountRepo().GetUserAccountByEmail(context.Background(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			// Don't reveal if email exists or not for security
			return nil
		}
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	// Generate password reset token (you can use JWT or random token)
	resetToken, err := s.generatePasswordResetToken(userAccount.Id)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// TODO: Store reset token in Redis with expiry (15 minutes)
	// TODO: Send email with reset link
	// For now, we'll just log it
	fmt.Printf("Password reset requested for user %s. Reset token: %s\n", userAccount.Email, resetToken)

	return nil
}

func (s *AuthService) ResetPassword(ctx *gin.Context, req dto.ResetPasswordRequest) error {
	// TODO: Validate reset token from Redis
	// For now, we'll just check if user exists
	userAccount, err := s.container.GetUserAccountRepo().GetUserAccountByEmail(context.Background(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to get user account: %w", err)
	}

	// TODO: Hash new password and update in database
	// For now, just return success
	fmt.Printf("Password reset for user %s\n", userAccount.Email)

	return nil
}

func (s *AuthService) ChangePassword(ctx *gin.Context, req dto.ChangePasswordRequest) error {
	// Get user account
	userAccount, err := s.container.GetUserAccountRepo().GetUserAccountByEmail(context.Background(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to get user account: %w", err)
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(userAccount.HashedPassword), []byte(req.CurrentPassword))
	if err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password in database
	err = s.container.GetUserAccountRepo().UpdateUserPassword(context.Background(), userAccount.Id, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// Helper function to generate password reset token
func (s *AuthService) generatePasswordResetToken(userID string) (string, error) {
	// Generate a simple reset token using JWT with short expiry
	tokenInput := &jwtService.GenerateTokenInput{
		UserId: userID,
		Email:  "", // Not needed for reset token
		Role:   "password_reset",
	}

	// For reset tokens, we can reuse the access token generation with shorter expiry
	resetToken, err := s.container.GetJWT().GenerateAccessToken(tokenInput)
	if err != nil {
		return "", err
	}

	return resetToken, nil
}
