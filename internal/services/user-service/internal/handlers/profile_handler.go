package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/toji-dev/go-shop/internal/pkg/constant"
	"github.com/toji-dev/go-shop/internal/pkg/response"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/services"
)

// ProfileHandler handles user profile-related requests
type ProfileHandler struct {
	userService *services.UserService
}

func NewProfileHandler(sc container.ServiceContainer) *ProfileHandler {
	return &ProfileHandler{
		userService: services.NewUserService(sc.GetUserProfileRepo()),
	}
}

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

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userIDRaw, exists := c.Get(constant.ContextKeyUserID)
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

func (h *ProfileHandler) GetProfileByID(c *gin.Context) {
	start := time.Now()

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
	currentUserID, exists := c.Get(constant.ContextKeyUserID)
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

		requestCount.WithLabelValues("/users/profile/:id", c.Request.Method).Inc()
		requestLatency.WithLabelValues("/users/profile/:id").Observe(time.Since(start).Seconds())
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

		requestCount.WithLabelValues("/users/profile/:id", c.Request.Method).Inc()
		requestLatency.WithLabelValues("/users/profile/:id").Observe(time.Since(start).Seconds())
		response.Success(c, "Public profile retrieved successfully", publicResponse)
	}
}

func (h *ProfileHandler) DeleteProfile(c *gin.Context) {
	userIDRaw, exists := c.Get(constant.ContextKeyUserID)
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
		constant.ContextKeyUserID: userIDStr,
		"status":                  "deleted",
	})
}

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goshop",
			Name:      "http_requests_total",
			Help:      "Tổng số request nhận được",
		},
		[]string{"endpoint", "method"},
	)

	requestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "goshop",
			Name:      "request_latency_seconds",
			Help:      "Thời gian xử lý request",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)
)

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestLatency)
}
