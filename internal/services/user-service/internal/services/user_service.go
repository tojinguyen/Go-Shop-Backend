package services

import (
	"context"
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

	userID := uuid.New()
	role := req.Role
	if role == "" {
		role = "user"
	}

	params := sqlc.CreateUserProfileParams{
		UserID:           uuidToPgtype(userID),
		Email:            req.Email,
		FullName:         stringToPgText(req.FirstName + " " + req.LastName),
		Birthday:         stringToPgDate(""), // Cập nhật nếu có trường birthday
		Phone:            stringToPgText(req.Phone),
		Role:             stringToPgText(role),
		BannedAt:         nullPgTime(),
		AvatarUrl:        stringToPgText(""), // Cập nhật nếu có trường avatar
		Gender:           stringToPgText(""), // Cập nhật nếu có trường gender
		DefaultAddressID: nullPgUUID(),
	}
	profile, err := repo.CreateUserProfile(context.Background(), params)
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
