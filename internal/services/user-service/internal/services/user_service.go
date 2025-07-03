package services

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/pkg/converter"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
)

type UserService struct {
	container *container.ServiceContainer
}

func NewUserService(container *container.ServiceContainer) *UserService {
	return &UserService{
		container: container,
	}
}

func (s *UserService) CreateProfile(ctx *gin.Context, req dto.CreateUserRequest) (domain.UserProfile, error) {
	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		return domain.UserProfile{}, fmt.Errorf("user ID not found in context")
	}
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		return domain.UserProfile{}, fmt.Errorf("user ID is not a string")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("invalid user ID format: %w", err)
	}

	roleRaw, exists := ctx.Get("user_role")
	if !exists {
		return domain.UserProfile{}, fmt.Errorf("user role not found in context")
	}
	role, ok := roleRaw.(string)
	if !ok {
		return domain.UserProfile{}, fmt.Errorf("user role is not a string")
	}

	params := sqlc.CreateUserProfileParams{
		UserID:           converter.UUIDToPgUUID(userID),
		Email:            req.Email,
		FullName:         req.FullName,
		Birthday:         converter.StringToPgDate(req.Birthday),
		Phone:            req.Phone,
		UserRole:         role,
		BannedAt:         converter.NullPgTime(),
		AvatarUrl:        req.AvatarURL,
		DefaultAddressID: converter.NullPgUUID(),
	}

	profile, err := s.container.GetUserProfileRepo().CreateUserProfile(ctx, params)
	if err != nil {
		return domain.UserProfile{}, err
	}
	return *profile, nil
}

func (s *UserService) GetProfile(ctx *gin.Context, userID string) (domain.UserProfile, error) {
	profile, err := s.container.GetUserProfileRepo().GetUserProfileByID(ctx, userID)
	if err != nil {
		return domain.UserProfile{}, err
	}
	return *profile, nil
}

func (s *UserService) UpdateProfile(ctx *gin.Context, req dto.UpdateUserRequest) (domain.UserProfile, error) {
	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		return domain.UserProfile{}, fmt.Errorf("user ID not found in context")
	}
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		return domain.UserProfile{}, fmt.Errorf("user ID is not a string")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Get current profile to preserve existing values
	currentProfile, err := s.container.GetUserProfileRepo().GetUserProfileByID(ctx, userIDStr)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("failed to get current profile: %w", err)
	}

	// Build update parameters
	params := sqlc.UpdateUserProfileParams{
		UserID:           converter.UUIDToPgUUID(userID),
		Email:            currentProfile.Email, // Keep existing email
		FullName:         currentProfile.FullName,
		Birthday:         converter.StringToPgDate(currentProfile.Birthday),
		Phone:            currentProfile.Phone,
		UserRole:         currentProfile.Role,
		BannedAt:         converter.StringToPgTime(currentProfile.BannedAt),
		AvatarUrl:        currentProfile.AvatarURL,
		Gender:           currentProfile.Gender,
		DefaultAddressID: converter.StringToPgUUID(currentProfile.DefaultAddressID),
	}

	// Update only provided fields
	if req.FullName != "" {
		params.FullName = req.FullName
	}

	if req.Phone != "" {
		params.Phone = req.Phone
	}

	if req.AvatarURL != "" {
		params.AvatarUrl = req.AvatarURL
	}

	if req.Birthday != "" {
		params.Birthday = converter.StringToPgDate(req.Birthday)
	}

	profile, err := s.container.GetUserProfileRepo().UpdateUserProfile(ctx, params)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("failed to update profile: %w", err)
	}
	return *profile, nil
}

func (s *UserService) DeleteProfile(ctx *gin.Context, userID string) error {
	// Verify that the user exists before attempting to delete
	_, err := s.container.GetUserProfileRepo().GetUserProfileByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Perform soft delete
	err = s.container.GetUserProfileRepo().DeleteProfile(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}
	return nil
}
