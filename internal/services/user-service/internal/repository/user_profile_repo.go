package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/domain"
)

type UserProfileRepository interface {
	CreateUserProfile(ctx context.Context, params sqlc.CreateUserProfileParams) (*domain.UserProfile, error)
}

type userProfileRepository struct {
	db      *postgresql_infra.PostgreSQLService
	queries *sqlc.Queries
}

func NewUserProfileRepository(db *postgresql_infra.PostgreSQLService) UserProfileRepository {
	return &userProfileRepository{
		db:      db,
		queries: sqlc.New(db.GetPool()),
	}
}

func (r *userProfileRepository) CreateUserProfile(ctx context.Context, params sqlc.CreateUserProfileParams) (*domain.UserProfile, error) {
	profile, err := r.queries.CreateUserProfile(ctx, params)
	if err != nil {
		return nil, err
	}
	return &domain.UserProfile{
		UserID:           profile.UserID.String(),
		Email:            profile.Email,
		FullName:         profile.FullName,
		Birthday:         pgDateToString(profile.Birthday),
		Phone:            profile.Phone,
		Role:             profile.Role,
		BannedAt:         pgTimeToString(profile.BannedAt),
		AvatarURL:        profile.AvatarUrl,
		Gender:           profile.Gender,
		DefaultAddressID: pgUUIDToString(profile.DefaultAddressID),
		CreatedAt:        pgTimeToString(profile.CreatedAt),
		UpdatedAt:        pgTimeToString(profile.UpdatedAt),
	}, nil
}

// Helper functions
func pgTimeToString(t pgtype.Timestamptz) string {
	if t.Valid {
		return t.Time.Format(time.RFC3339)
	}
	return ""
}

func pgDateToString(d pgtype.Date) string {
	if d.Valid {
		return d.Time.Format("2006-01-02")
	}
	return ""
}

func pgUUIDToString(u pgtype.UUID) string {
	if u.Valid {
		return u.String()
	}
	return ""
}
