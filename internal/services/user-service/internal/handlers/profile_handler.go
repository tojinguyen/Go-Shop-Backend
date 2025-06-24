package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/your-username/go-shop/internal/services/user-service/internal/config"
	"github.com/your-username/go-shop/internal/services/user-service/internal/dto"
	jwtService "github.com/your-username/go-shop/internal/services/user-service/internal/pkg/jwt"
	"github.com/your-username/go-shop/internal/services/user-service/internal/pkg/response"
)

type ProfileHandler struct {
	jwtService jwtService.JwtService
	config     *config.Config
}

// NewAuthHandler creates a new auth handler
func NewProfileHandler(jwtSvc jwtService.JwtService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		jwtService: jwtSvc,
		config:     cfg,
	}
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
	user := &dto.UserInfo{
		ID:        userID.(string),
		Email:     "user@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Role:      "user",
	}

	response.Success(c, "Profile retrieved successfully", user)
}
