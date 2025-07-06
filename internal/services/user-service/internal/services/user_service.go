package services

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/pkg/converter"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/repository"
)

type UserService struct {
	userProfileRepo repository.UserProfileRepository
}

func NewUserService(userProfileRepo repository.UserProfileRepository) *UserService {
	return &UserService{
		userProfileRepo: userProfileRepo,
	}
}

func (s *UserService) CreateProfile(ctx *gin.Context, req dto.CreateUserProfileRequest) (domain.UserProfile, error) {
	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		log.Printf("user ID not found in context")
		return domain.UserProfile{}, fmt.Errorf("user ID not found in context")
	}
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		log.Printf("user ID is not a string")
		return domain.UserProfile{}, fmt.Errorf("user ID is not a string")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("invalid user ID format: %v", err)
		return domain.UserProfile{}, fmt.Errorf("invalid user ID format: %w", err)
	}

	roleRaw, exists := ctx.Get("user_role")
	if !exists {
		log.Printf("user role not found in context")
		return domain.UserProfile{}, fmt.Errorf("user role not found in context")
	}
	role, ok := roleRaw.(string)
	if !ok {
		log.Printf("user role is not a string")
		return domain.UserProfile{}, fmt.Errorf("user role is not a string")
	}

	email, exists := ctx.Get("user_email")
	if !exists {
		log.Printf("email not found in context")
		return domain.UserProfile{}, fmt.Errorf("email not found in context: %w", err)
	}
	emailStr, ok := email.(string)
	if !ok {
		log.Printf("email is not a string")
		return domain.UserProfile{}, fmt.Errorf("email is not a string")
	}

	params := sqlc.CreateUserProfileParams{
		UserID:           converter.UUIDToPgUUID(userID),
		Email:            emailStr,
		FullName:         req.FullName,
		Birthday:         converter.StringToPgDate(req.Birthday),
		Phone:            req.Phone,
		UserRole:         role,
		BannedAt:         converter.NullPgTime(),
		AvatarUrl:        req.AvatarURL,
		DefaultAddressID: converter.NullPgUUID(),
		Gender:           req.Gender,
	}

	profile, err := s.userProfileRepo.CreateUserProfile(ctx, params)
	if err != nil {
		log.Printf("failed to create user profile: %v", err)
		return domain.UserProfile{}, err
	}
	return *profile, nil
}

func (s *UserService) GetProfile(ctx *gin.Context, userID string) (domain.UserProfile, error) {
	profile, err := s.userProfileRepo.GetUserProfileByID(ctx, userID)
	if err != nil {
		return domain.UserProfile{}, err
	}
	return *profile, nil
}

func (s *UserService) UpdateProfile(ctx *gin.Context, req dto.UpdateUserProfileRequest) (domain.UserProfile, error) {
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
	currentProfile, err := s.userProfileRepo.GetUserProfileByID(ctx, userIDStr)
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

	profile, err := s.userProfileRepo.UpdateUserProfile(ctx, params)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("failed to update profile: %w", err)
	}
	return *profile, nil
}

func (s *UserService) DeleteProfile(ctx *gin.Context, userID string) error {
	// Verify that the user exists before attempting to delete
	_, err := s.userProfileRepo.GetUserProfileByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Perform soft delete
	err = s.userProfileRepo.DeleteProfile(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}
	return nil
}

func (s *UserService) GetProfileByID(ctx *gin.Context, userID string) (domain.UserProfile, error) {
	// Validate UUID format
	_, err := uuid.Parse(userID)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("invalid user ID format")
	}

	// Get user profile by ID
	profile, err := s.userProfileRepo.GetUserProfileByID(ctx, userID)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("user not found")
	}
	return *profile, nil
}
