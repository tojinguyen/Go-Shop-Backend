package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/your-username/go-shop/internal/services/user-service/internal/config"
	errorConstants "github.com/your-username/go-shop/internal/services/user-service/internal/pkg/errors"
	jwtService "github.com/your-username/go-shop/internal/services/user-service/internal/pkg/kwt"
	"github.com/your-username/go-shop/internal/services/user-service/internal/pkg/response"
	"github.com/your-username/go-shop/internal/services/user-service/internal/pkg/validation"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	jwtService jwtService.JwtService
	config     *config.Config
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(jwtSvc jwtService.JwtService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		jwtService: jwtSvc,
		config:     cfg,
	}
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Username        string `json:"username"`
	Phone           string `json:"phone"`
}

// RefreshTokenRequest represents the refresh token request payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	User         *UserInfo `json:"user"`
}

// UserInfo represents user information in auth response
type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Role      string `json:"role"`
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
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
	mockUser := &UserInfo{
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

	authResponse := &AuthResponse{
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
	var req RegisterRequest
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

	// Validate names
	if isValid, nameErrors := validation.ValidateName(req.FirstName); !isValid {
		response.BadRequest(c, "INVALID_FIRST_NAME", "Invalid first name", strings.Join(nameErrors, "; "))
		return
	}

	if isValid, nameErrors := validation.ValidateName(req.LastName); !isValid {
		response.BadRequest(c, "INVALID_LAST_NAME", "Invalid last name", strings.Join(nameErrors, "; "))
		return
	}

	// Validate username if provided
	if req.Username != "" {
		if isValid, usernameErrors := validation.ValidateUsername(req.Username); !isValid {
			response.BadRequest(c, "INVALID_USERNAME", "Invalid username", strings.Join(usernameErrors, "; "))
			return
		}
	}

	// Validate phone if provided
	if req.Phone != "" && !validation.ValidatePhoneNumber(req.Phone) {
		response.BadRequest(c, "INVALID_PHONE", "Invalid phone number format", "")
		return
	}

	// Here you would typically:
	// 1. Check if user already exists
	// 2. Hash the password
	// 3. Save user to database
	// 4. Send verification email

	// For demo purposes, we'll create a mock user
	newUser := &UserInfo{
		ID:        "user-456", // This would be generated by database
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Role:      "user",
	}

	response.Created(c, "User registered successfully", newUser)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
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

	authResponse := &AuthResponse{
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

// GetProfile returns the current user's profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_AUTHENTICATED", "User not authenticated")
		return
	}

	// Here you would typically fetch user details from database
	// For demo purposes, we'll return mock data
	user := &UserInfo{
		ID:        userID.(string),
		Email:     "user@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Role:      "user",
	}

	response.Success(c, "Profile retrieved successfully", user)
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
	tokenInfo := gin.H{
		"valid":      true,
		"user_id":    claims.UserId,
		"user_email": claims.Email,
		"user_role":  claims.Role,
		"expires_at": claims.ExpiresAt,
		"issued_at":  claims.IssuedAt,
		"issuer":     claims.Issuer,
	}

	response.Success(c, "Token is valid", tokenInfo)
}
