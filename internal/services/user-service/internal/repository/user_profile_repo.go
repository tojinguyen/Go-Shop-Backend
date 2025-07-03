package repository

import (
	"context"

	"github.com/toji-dev/go-shop/internal/pkg/converter"
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
		Birthday:         converter.PgDateToString(profile.Birthday),
		Phone:            profile.Phone,
		Role:             profile.UserRole,
		BannedAt:         converter.PgTimeToString(profile.BannedAt),
		AvatarURL:        profile.AvatarUrl,
		Gender:           profile.Gender,
		DefaultAddressID: converter.PgUUIDToString(profile.DefaultAddressID),
		CreatedAt:        converter.PgTimeToString(profile.CreatedAt),
		UpdatedAt:        converter.PgTimeToString(profile.UpdatedAt),
	}, nil
}
