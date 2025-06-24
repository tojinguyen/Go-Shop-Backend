package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/your-username/go-shop/internal/services/user-service/internal/config"
	"github.com/your-username/go-shop/internal/services/user-service/internal/dto"
	jwtService "github.com/your-username/go-shop/internal/services/user-service/internal/pkg/jwt"
	"github.com/your-username/go-shop/internal/services/user-service/internal/pkg/response"
)

type ProfileHandler struct {
	config     *config.Config
	jwtService jwtService.JwtService
}

// NewAuthHandler creates a new auth handler
func NewProfileHandler(jwtSvc jwtService.JwtService, cfg *config.Config) *ProfileHandler {
	return &ProfileHandler{
		jwtService: jwtSvc,
		config:     cfg,
	}
}

// GetProfile returns the current user's profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
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

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_AUTHENTICATED", "User not authenticated")
		return
	}
	//TODO: update profile

	response.Success(c, "Profile retrieved successfully", userID)
}

func (h *ProfileHandler) GetProfileByID(c *gin.Context) {
	// Get user ID from URL parameters
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "INVALID_USER_ID", "User ID is required", "User not authenticated")
		return
	}

	// Here you would typically fetch user details from database
	// For demo purposes, we'll return mock data
	user := &dto.UserInfo{
		ID:        userID,
		Email:     "user@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Role:      "user",
	}
	response.Success(c, "Profile retrieved successfully", user)
}

func (h *ProfileHandler) DeleteProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_AUTHENTICATED", "User not authenticated")
		return
	}
	// Here you would typically delete the user from the database
	// For demo purposes, we'll just return a success message
	response.Success(c, "Profile deleted successfully", userID)
}
