package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/response"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/validation"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/services"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(sc container.ServiceContainer) *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(&sc),
	}
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

	res, err := h.authService.Register(c, req)

	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			response.Conflict(c, "USER_ALREADY_EXISTS", "User with this email already exists")
			return
		}
		response.InternalServerError(c, "REGISTRATION_FAILED", "Failed to register user")
		return
	}

	response.Created(c, "User registered successfully", res)
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

	// Use AuthService to handle login
	authResponse, err := h.authService.Login(c, req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid email or password") {
			response.Unauthorized(c, "INVALID_CREDENTIALS", "Invalid email or password")
			return
		}
		response.InternalServerError(c, "LOGIN_FAILED", "Failed to login user")
		return
	}

	response.Success(c, "Login successful", authResponse)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request payload", err.Error())
		return
	}

	// Use AuthService to handle token refresh
	authResponse, err := h.authService.RefreshToken(c, req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid refresh token") {
			response.Unauthorized(c, "REFRESH_TOKEN_INVALID", "Invalid refresh token")
			return
		}
		response.InternalServerError(c, "TOKEN_REFRESH_FAILED", "Failed to refresh token")
		return
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
	tokenInfo, err := h.authService.ValidateToken(c, token)
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			response.Unauthorized(c, "TOKEN_EXPIRED", "Token has expired")
			return
		}
		if strings.Contains(err.Error(), "invalid") {
			response.Unauthorized(c, "TOKEN_INVALID", "Invalid token")
			return
		}
		response.Unauthorized(c, "TOKEN_VALIDATION_FAILED", "Token validation failed")
		return
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
