package services

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
	// Lấy repository từ ServiceContainer
	repo := s.container.GetUserProfileRepo()

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
		UserID:           uuidToPgtype(userID),
		Email:            req.Email,
		FullName:         req.FullName,
		Birthday:         stringToPgDate(req.Birthday),
		Phone:            req.Phone,
		Role:             role,
		BannedAt:         nullPgTime(),
		AvatarUrl:        req.AvatarURL,
		DefaultAddressID: nullPgUUID(),
	}

	profile, err := repo.CreateUserProfile(ctx, params)
	if err != nil {
		return domain.UserProfile{}, err
	}
	return *profile, nil
}

// Helper functions for pgtype
func uuidToPgtype(id uuid.UUID) pgtype.UUID {
	var u pgtype.UUID
	u.Scan(id.String())
	return u
}
func stringToPgText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}
func stringToPgDate(s string) pgtype.Date {
	var d pgtype.Date
	if s != "" {
		d.Time, _ = time.Parse("2006-01-02", s)
		d.Valid = true
	}
	return d
}
func nullPgTime() pgtype.Timestamptz {
	return pgtype.Timestamptz{}
}
func nullPgUUID() pgtype.UUID {
	return pgtype.UUID{}
}
