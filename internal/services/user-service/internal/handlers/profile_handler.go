package handlers

import (
	"github.com/gin-gonic/gin"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	redis_infra "github.com/toji-dev/go-shop/internal/pkg/infra/redis-infra"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	jwtService "github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/jwt"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/response"
)

type ProfileHandler struct {
	config       *config.Config
	jwtService   jwtService.JwtService
	pgService    *postgresql_infra.PostgreSQLService
	redisService *redis_infra.RedisService
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(jwtSvc jwtService.JwtService, cfg *config.Config, pgSvc *postgresql_infra.PostgreSQLService, redisSvc *redis_infra.RedisService) *ProfileHandler {
	return &ProfileHandler{
		jwtService:   jwtSvc,
		config:       cfg,
		pgService:    pgSvc,
		redisService: redisSvc,
	}
}

// GetProfile returns the current user's profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// Here you would typically fetch user details from database
	// For demo purposes, we'll return mock data
	user := &dto.UserInfo{}

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
	user := &dto.UserInfo{}
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
