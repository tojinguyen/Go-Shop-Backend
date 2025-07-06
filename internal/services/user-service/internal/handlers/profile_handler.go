package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/response"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/services"
)

// ProfileHandler handles user profile-related requests
type ProfileHandler struct {
	userService *services.UserService
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(sc container.ServiceContainer) *ProfileHandler {
	return &ProfileHandler{
		userService: services.NewUserService(sc.GetUserProfileRepo()),
	}
}

// CreateProfile creates a new user profile
// @Summary Create user profile
// @Description Create a new user profile with detailed information
// @Tags profile
// @Accept json
// @Produce json
// @Param request body dto.CreateUserProfileRequest true "Create profile request"
// @Success 200 {object} map[string]interface{} "Profile created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 409 {object} map[string]interface{} "User already exists"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /profile [post]
func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	// Bind the request body to CreateUserRequest
	var req dto.CreateUserProfileRequest
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
// @Summary Get current user profile
// @Description Get the authenticated user's profile information
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Profile retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /profile [get]
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

// UpdateProfile updates the current user's profile
// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateUserProfileRequest true "Update profile request"
// @Success 200 {object} map[string]interface{} "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /profile [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	// Bind the request body to UpdateUserRequest
	var req dto.UpdateUserProfileRequest
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

// GetProfileByID gets a user profile by ID
// @Summary Get user profile by ID
// @Description Get a user's profile by their ID (public or private view based on authentication)
// @Tags profile
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{} "Profile retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /profile/{id} [get]
func (h *ProfileHandler) GetProfileByID(c *gin.Context) {
	// Get user ID from URL parameters
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "INVALID_USER_ID", "User ID is required", "User ID parameter is missing")
		return
	}

	// Get user profile by ID
	userProfile, err := h.userService.GetProfileByID(c, userID)
	if err != nil {
		if err.Error() == "user not found" {
			response.NotFound(c, "USER_NOT_FOUND", "User profile not found")
			return
		}
		if err.Error() == "invalid user ID format" {
			response.BadRequest(c, "INVALID_USER_ID", "Invalid user ID format", "User ID must be a valid UUID")
			return
		}
		response.InternalServerError(c, "PROFILE_RETRIEVAL_FAILED", "Failed to retrieve profile")
		return
	}

	// Check if the requesting user is viewing their own profile
	currentUserID, exists := c.Get("user_id")
	isOwnProfile := exists && currentUserID == userID

	if isOwnProfile {
		// Return full profile information for own profile
		userResponse := &dto.UserResponse{
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
		response.Success(c, "Profile retrieved successfully", userResponse)
	} else {
		// Return limited public profile information
		publicResponse := &dto.PublicUserResponse{
			ID:        userProfile.UserID,
			FullName:  userProfile.FullName,
			Avatar:    userProfile.AvatarURL,
			Role:      userProfile.Role,
			CreatedAt: userProfile.CreatedAt,
		}
		response.Success(c, "Public profile retrieved successfully", publicResponse)
	}
}

// DeleteProfile deletes the current user's profile
// @Summary Delete user profile
// @Description Delete the authenticated user's profile permanently
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Profile deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /profile [delete]
func (h *ProfileHandler) DeleteProfile(c *gin.Context) {
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
