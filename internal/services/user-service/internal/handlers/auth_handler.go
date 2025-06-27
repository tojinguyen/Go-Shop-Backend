package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	redis_infra "github.com/toji-dev/go-shop/internal/pkg/infra/redis-infra"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	errorConstants "github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/errors"
	jwtService "github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/jwt"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/response"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/validation"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	jwtService   jwtService.JwtService
	config       *config.Config
	pgService    *postgresql_infra.PostgreSQLService
	redisService *redis_infra.RedisService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(jwtSvc jwtService.JwtService, cfg *config.Config, pgSvc *postgresql_infra.PostgreSQLService, redisSvc *redis_infra.RedisService) *AuthHandler {
	return &AuthHandler{
		jwtService:   jwtSvc,
		config:       cfg,
		pgService:    pgSvc,
		redisService: redisSvc,
	}
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request payload", err.Error())
		return
	}

	// Validate email format
	if !validation.ValidateEmail(req.Email) {
		response.BadRequest(c, "INVALID_EMAIL", "Invalid email format", "")
		return
	}

	// Here you would typically:
	// 1. Find user by email in database
	// 2. Verify password hash
	// 3. Check if user is active

	// For demo purposes, we'll create a mock user
	mockUser := &dto.UserInfo{
		ID:        "user-123",
		Email:     req.Email,
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Role:      "user",
	}

	// Generate tokens
	tokenInput := &jwtService.GenerateTokenInput{
		UserId: mockUser.ID,
		Email:  mockUser.Email,
		Role:   mockUser.Role,
	}

	accessToken, err := h.jwtService.GenerateAccessToken(tokenInput)
	if err != nil {
		response.InternalServerError(c, "TOKEN_GENERATION_FAILED", "Failed to generate access token")
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(tokenInput)
	if err != nil {
		response.InternalServerError(c, "TOKEN_GENERATION_FAILED", "Failed to generate refresh token")
		return
	}

	authResponse := &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.config.JWT.AccessTokenTTL.Seconds()),
		User:         mockUser,
	}

	response.Success(c, "Login successful", authResponse)
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request payload", err.Error())
		return
	}

	// Validate email
	if !validation.ValidateEmail(req.Email) {
		response.BadRequest(c, "INVALID_EMAIL", "Invalid email format", "")
		return
	}

	// Validate password
	isValidPassword, passwordErrors := validation.ValidatePassword(req.Password)
	if !isValidPassword {
		response.BadRequest(c, "WEAK_PASSWORD", "Password does not meet requirements", strings.Join(passwordErrors, "; "))
		return
	}

	// Check password confirmation
	if req.Password != req.ConfirmPassword {
		response.BadRequest(c, "PASSWORD_MISMATCH", "Passwords do not match", "")
		return
	}

	response.Created(c, "User registered successfully", nil)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request payload", err.Error())
		return
	}

	// Validate refresh token
	claims, err := h.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		var errorCode string
		switch err {
		case errorConstants.ErrTokenExpired:
			errorCode = "REFRESH_TOKEN_EXPIRED"
		case errorConstants.ErrTokenInvalid:
			errorCode = "REFRESH_TOKEN_INVALID"
		default:
			errorCode = "REFRESH_TOKEN_VALIDATION_FAILED"
		}

		response.Unauthorized(c, errorCode, err.Error())
		return
	}

	// Generate new access token
	tokenInput := &jwtService.GenerateTokenInput{
		UserId: claims.UserId,
		Email:  claims.Email,
		Role:   claims.Role,
	}

	newAccessToken, err := h.jwtService.GenerateAccessToken(tokenInput)
	if err != nil {
		response.InternalServerError(c, "TOKEN_GENERATION_FAILED", "Failed to generate access token")
		return
	}

	// Optionally generate new refresh token
	newRefreshToken, err := h.jwtService.GenerateRefreshToken(tokenInput)
	if err != nil {
		response.InternalServerError(c, "TOKEN_GENERATION_FAILED", "Failed to generate refresh token")
		return
	}

	authResponse := &dto.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.config.JWT.AccessTokenTTL.Seconds()),
	}

	response.Success(c, "Token refreshed successfully", authResponse)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Here you would typically:
	// 1. Blacklist the current token
	// 2. Clear any user sessions
	// 3. Log the logout event

	response.Success(c, "Logout successful", nil)
}

// ValidateToken validates the provided token
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	// This endpoint is useful for other services to validate tokens
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.BadRequest(c, "MISSING_TOKEN", "Authorization header is required", "")
		return
	}

	// Extract token
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		response.BadRequest(c, "INVALID_TOKEN_FORMAT", "Invalid authorization header format", "")
		return
	}

	// Validate token
	claims, err := h.jwtService.ValidateAccessToken(c.Request.Context(), token)
	if err != nil {
		var errorCode string
		switch err {
		case errorConstants.ErrTokenExpired:
			errorCode = "TOKEN_EXPIRED"
		case errorConstants.ErrTokenInvalid:
			errorCode = "TOKEN_INVALID"
		default:
			errorCode = "TOKEN_VALIDATION_FAILED"
		}

		response.Unauthorized(c, errorCode, err.Error())
		return
	}

	// Return token information
	tokenInfo := &dto.TokenValidationResponse{
		Valid:     true,
		UserID:    claims.UserId,
		UserEmail: claims.Email,
		UserRole:  claims.Role,
		ExpiresAt: claims.ExpiresAt,
		IssuedAt:  claims.IssuedAt,
		Issuer:    claims.Issuer,
	}

	response.Success(c, "Token is valid", tokenInfo)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request payload", err.Error())
		return
	}

	// Validate email format
	if !validation.ValidateEmail(req.Email) {
		response.BadRequest(c, "INVALID_EMAIL", "Invalid email format", "")
		return
	}

	// Here you would typically:
	// 1. Check if user exists by email
	// 2. Generate a password reset token
	// 3. Send reset link via email

	response.Success(c, "Password reset link sent successfully", nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request payload", err.Error())
		return
	}

	// Validate email format
	if !validation.ValidateEmail(req.Email) {
		response.BadRequest(c, "INVALID_EMAIL", "Invalid email format", "")
		return
	}

	response.Success(c, "Password reset successfully", nil)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request payload", err.Error())
		return
	}

	// Validate current password
	if req.CurrentPassword == "" {
		response.BadRequest(c, "MISSING_CURRENT_PASSWORD", "Current password is required", "")
		return
	}

	// Validate new password
	isValidPassword, passwordErrors := validation.ValidatePassword(req.NewPassword)
	if !isValidPassword {
		response.BadRequest(c, "WEAK_PASSWORD", "New password does not meet requirements", strings.Join(passwordErrors, "; "))
		return
	}

	// Check if new password matches confirmation
	if req.NewPassword != req.ConfirmPassword {
		response.BadRequest(c, "PASSWORD_MISMATCH", "New passwords do not match", "")
		return
	}

	// Here you would typically:
	// 1. Verify current password against stored hash
	// 2. Hash the new password
	// 3. Update user password in database

	response.Success(c, "Password changed successfully", nil)
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request payload", err.Error())
		return
	}

	// Validate OTP
	var valid, otpErrors = validation.ValidateOTP(req.OTP)
	if !valid {
		response.BadRequest(c, "INVALID_OTP", "OTP is invalid", otpErrors)
		return
	}

	// Here you would typically:
	// 1. Verify the OTP against the stored value
	// 2. Mark the user as verified if successful

	response.Success(c, "OTP verified successfully", nil)
}
