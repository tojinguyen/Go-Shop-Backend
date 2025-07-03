package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/response"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/services"
)

type ProfileHandler struct {
	userService *services.UserService
}

// NewAuthHandler creates a new auth handler
func NewProfileHandler(sc container.ServiceContainer) *ProfileHandler {
	return &ProfileHandler{
		userService: services.NewUserService(&sc),
	}
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	// Bind the request body to CreateUserRequest
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request payload", err.Error())
		return
	}

	userProfile, err := h.userService.CreateProfile(c, req)

	if err != nil {
		if err.Error() == "user already exists" {
			response.Conflict(c, "USER_ALREADY_EXISTS", "User with this email already exists")
			return
		}
		response.InternalServerError(c, "PROFILE_CREATION_FAILED", "Failed to create profile")
		return
	}

	userProfileResponse := &dto.UserResponse{
		ID:               userProfile.UserID,
		Email:            userProfile.Email,
		FullName:         userProfile.FullName,
		Birthday:         userProfile.Birthday,
		Phone:            userProfile.Phone,
		Avatar:           userProfile.AvatarURL,
		Role:             userProfile.Role,
		Gender:           userProfile.Gender,
		DefaultAddressID: userProfile.DefaultAddressID,
		CreatedAt:        userProfile.CreatedAt,
		UpdatedAt:        userProfile.UpdatedAt,
	}

	response.Success(c, "Profile created successfully", userProfileResponse)
}

// GetProfile returns the current user's profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_AUTHENTICATED", "User not authenticated")
		return
	}
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		response.BadRequest(c, "INVALID_USER_ID", "User ID is not a valid string", "User not authenticated")
		return
	}

	userProfile, err := h.userService.GetProfile(c, userIDStr)

	if err != nil {
		response.InternalServerError(c, "PROFILE_RETRIEVAL_FAILED", "Failed to retrieve profile")
		return
	}

	response.Success(c, "Profile retrieved successfully", userProfile)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	// Bind the request body to UpdateUserRequest
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request payload", err.Error())
		return
	}

	// Update the user profile
	updatedProfile, err := h.userService.UpdateProfile(c, req)
	if err != nil {
		if err.Error() == "user not found" {
			response.NotFound(c, "USER_NOT_FOUND", "User profile not found")
			return
		}
		response.InternalServerError(c, "PROFILE_UPDATE_FAILED", "Failed to update profile")
		return
	}

	response.Success(c, "Profile retrieved successfully", updatedProfile)
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
	user := &dto.UserInfo{}
	response.Success(c, "Profile retrieved successfully", user)
}

func (h *ProfileHandler) DeleteProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_AUTHENTICATED", "User not authenticated")
		return
	}
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		response.BadRequest(c, "INVALID_USER_ID", "User ID is not a valid string", "User not authenticated")
		return
	}

	// Delete the user profile
	err := h.userService.DeleteProfile(c, userIDStr)
	if err != nil {
		if err.Error() == "user not found" {
			response.NotFound(c, "USER_NOT_FOUND", "User profile not found")
			return
		}
		response.InternalServerError(c, "PROFILE_DELETION_FAILED", "Failed to delete profile")
		return
	}

	response.Success(c, "Profile deleted successfully", map[string]string{
		"user_id": userIDStr,
		"status":  "deleted",
	})
}
