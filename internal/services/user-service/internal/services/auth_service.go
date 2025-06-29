package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	jwtService "github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/jwt"
	timeUtils "github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/time"
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
		log.Println("Error checking user existence:", err)
		return dto.RegisterResponse{}, fmt.Errorf("failed to check user existence: %w", err)
	}

	if exists {
		log.Printf("User with email %s already exists", req.Email)
		return dto.RegisterResponse{}, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
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
		log.Println("Error creating user account:", err)
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

	// Add token to Redis blacklist
	blacklistKey := fmt.Sprintf("blacklisted_token:%s", token)

	// Calculate expiry time based on token's remaining lifetime
	secondsUntilExpiry := timeUtils.GetSecondsUntilExpiry(claims.ExpiresAt)

	if secondsUntilExpiry > 0 {
		// Token is still valid, blacklist it until it naturally expires
		remainingTTL := timeUtils.CalculateDurationFromSeconds(secondsUntilExpiry)

		// Store token in Redis blacklist with remaining TTL
		err = s.container.GetRedis().Set(blacklistKey, "blacklisted", remainingTTL)
		if err != nil {
			// Log error but don't fail the logout
			fmt.Printf("Warning: Failed to blacklist token in Redis: %v\n", err)
		} else {
			fmt.Printf("Token successfully blacklisted in Redis with TTL: %v\n", remainingTTL)
		}
	}

	// Log the logout
	fmt.Printf("User %s logged out successfully\n", claims.UserId)

	return nil
}

func (s *AuthService) ForgotPassword(ctx *gin.Context, req dto.ForgotPasswordRequest) error {
	// Check if user exists
	userAccount, err := s.container.GetUserAccountRepo().GetUserAccountByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			// Don't reveal if email exists or not for security
			return nil
		}
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	// Generate password reset token
	resetToken, err := s.generatePasswordResetToken(userAccount.Id)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Store reset token in Redis with 15 minutes expiry
	resetTokenKey := fmt.Sprintf("password_reset_token:%s", userAccount.Id)
	err = s.container.GetRedis().Set(resetTokenKey, resetToken, 15*time.Minute)
	if err != nil {
		return fmt.Errorf("failed to store reset token: %w", err)
	}

	// Get frontend URL from config
	frontendURL := s.container.GetConfig().App.FrontendURL
	if frontendURL == "" {
		frontendURL = "http://localhost:3000" // fallback
	}

	// Prepare email data
	emailData := struct {
		UserName   string
		ResetToken string
		ResetLink  string
	}{
		UserName:   userAccount.Email, // You might want to add a name field to user account
		ResetToken: resetToken,
		ResetLink:  fmt.Sprintf("%s/reset-password?token=%s", frontendURL, resetToken),
	}

	// Send password reset email using template
	subject := "Password Reset Request - Go-Shop"

	// Try to use template first, fallback to HTML if template fails
	err = s.container.GetEmail().SendTemplateEmail([]string{userAccount.Email}, subject, "password_reset.html", emailData)
	if err != nil {
		// Fallback to simple HTML email if template fails
		htmlBody := fmt.Sprintf(`
			<html>
			<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
				<div style="background-color: #f9f9f9; padding: 30px; border-radius: 10px; border: 1px solid #ddd;">
					<h2 style="color: #e74c3c; text-align: center;">Password Reset Request</h2>
					<p>Hello,</p>
					<p>You have requested to reset your password for your Go-Shop account.</p>
					<p>Please click the link below to reset your password:</p>
					<div style="text-align: center; margin: 20px 0;">
						<a href="%s" style="background-color: #3498db; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; font-weight: bold;">Reset Password</a>
					</div>
					<div style="background-color: #fff3cd; border: 1px solid #ffeaa7; color: #856404; padding: 15px; border-radius: 5px; margin: 20px 0;">
						<strong>Security Notice:</strong> This link will expire in 15 minutes.
					</div>
					<p>If you did not request this password reset, please ignore this email.</p>
					<br>
					<p>Best regards,<br>Go-Shop Team</p>
				</div>
			</body>
			</html>
		`, emailData.ResetLink)

		err = s.container.GetEmail().SendHTMLEmail([]string{userAccount.Email}, subject, htmlBody)
		if err != nil {
			// Don't fail the request if email fails, but log the error
			fmt.Printf("Warning: Failed to send password reset email to %s: %v\n", userAccount.Email, err)
			return nil
		}
	}

	fmt.Printf("Password reset email sent successfully to %s\n", userAccount.Email)

	return nil
}

func (s *AuthService) ResetPassword(ctx *gin.Context, req dto.ResetPasswordRequest) error {
	// Validate the reset token
	claims, err := s.container.GetJWT().ValidateAccessToken(context.Background(), req.Token)
	if err != nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	// Check if this is a password reset token
	if claims.Role != "password_reset" {
		return fmt.Errorf("invalid reset token")
	}

	// Check if the token exists in Redis
	resetTokenKey := fmt.Sprintf("password_reset_token:%s", claims.UserId)
	exists, err := s.container.GetRedis().Exists(resetTokenKey)
	if err != nil {
		return fmt.Errorf("failed to verify reset token: %w", err)
	}

	if !exists {
		return fmt.Errorf("reset token has expired or already been used")
	}

	// Get user account to ensure user still exists
	userAccount, err := s.container.GetUserAccountRepo().GetUserAccountByID(context.Background(), claims.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to get user account: %w", err)
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

	// Remove the reset token from Redis to prevent reuse
	err = s.container.GetRedis().Delete(resetTokenKey)
	if err != nil {
		// Log error but don't fail the password reset
		fmt.Printf("Warning: Failed to remove reset token from Redis: %v\n", err)
	}

	fmt.Printf("Password reset successfully for user %s\n", userAccount.Email)
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

// IsTokenBlacklisted checks if a token is in the Redis blacklist
func (s *AuthService) IsTokenBlacklisted(token string) bool {
	blacklistKey := fmt.Sprintf("blacklisted_token:%s", token)

	exists, err := s.container.GetRedis().Exists(blacklistKey)
	if err != nil {
		// If Redis is unavailable, log error but don't block authentication
		fmt.Printf("Warning: Failed to check token blacklist: %v\n", err)
		return false
	}

	return exists
}

// Helper function to generate password reset token
func (s *AuthService) generatePasswordResetToken(userID string) (string, error) {
	// Generate a password reset token using JWT with 15 minutes expiry
	tokenInput := &jwtService.GenerateTokenInput{
		UserId: userID,
		Email:  "", // Not needed for reset token
		Role:   "password_reset",
	}

	// Generate a short-lived reset token
	resetToken, err := s.container.GetJWT().GenerateAccessToken(tokenInput)
	if err != nil {
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}

	return resetToken, nil
}
