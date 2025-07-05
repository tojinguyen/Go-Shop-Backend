package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/toji-dev/go-shop/internal/pkg/converter"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/domain"
)

// PostgresShopRepository implements the ShopRepository interface using PostgreSQL
type PostgresShopRepository struct {
	db      *postgresql_infra.PostgreSQLService
	queries *sqlc.Queries
}

// NewPostgresShopRepository creates a new PostgreSQL shop repository
func NewPostgresShopRepository(db *postgresql_infra.PostgreSQLService) ShopRepository {
	return &PostgresShopRepository{
		db:      db,
		queries: sqlc.New(db.GetPool()),
	}
}

// Create creates a new shop
func (r *PostgresShopRepository) Create(ctx context.Context, shop *domain.Shop) error {
	rating := converter.Float64ToPgNumeric(shop.Rating)

	sqlcParams := sqlc.CreateShopParams{
		ID:              converter.UUIDToPgUUID(shop.ID),
		OwnerID:         converter.UUIDToPgUUID(shop.OwnerID),
		ShopName:        shop.ShopName,
		AvatarUrl:       shop.AvatarURL,
		BannerUrl:       shop.BannerURL,
		ShopDescription: converter.StringToPgText(shop.ShopDescription),
		AddressID:       converter.UUIDToPgUUID(shop.AddressID),
		Phone:           shop.Phone,
		Email:           shop.Email,
		Rating:          rating,
		ActiveAt:        converter.TimePtrToPgTime(shop.ActiveAt),
	}

	createdShop, err := r.queries.CreateShop(ctx, sqlcParams)
	if err != nil {
		log.Println("Error creating shop:", err)
		return fmt.Errorf("failed to create shop: %w", err)
	}

	// Update the shop with the created values (like timestamps)
	if createdShop.CreatedAt.Valid {
		shop.CreatedAt = createdShop.CreatedAt.Time
	}
	if createdShop.UpdatedAt.Valid {
		shop.UpdatedAt = createdShop.UpdatedAt.Time
	}

	return nil
}

func (r *PostgresShopRepository) GetShopByID(ctx context.Context, shopID string) (*domain.Shop, error) {
	shopIDUUID := converter.StringToPgUUID(shopID)

	shop, err := r.queries.GetShopByID(ctx, shopIDUUID)
	if err != nil {
		log.Println("Error fetching shop by ID:", err)
		return nil, fmt.Errorf("failed to get shop by ID: %w", err)
	}

	return &domain.Shop{
		ID:              converter.PgUUIDToUUID(shop.ID),
		OwnerID:         converter.PgUUIDToUUID(shop.OwnerID),
		ShopName:        shop.ShopName,
		AvatarURL:       shop.AvatarUrl,
		BannerURL:       shop.BannerUrl,
		ShopDescription: converter.PgTextToStringPtr(shop.ShopDescription),
		AddressID:       converter.PgUUIDToUUID(shop.AddressID),
		Phone:           shop.Phone,
		Email:           shop.Email,
		Rating:          *converter.PgNumericToFloat64Ptr(shop.Rating),
		ActiveAt:        converter.PgTimeToTimePtr(shop.ActiveAt),
		BannedAt:        converter.PgTimeToTimePtr(shop.BannedAt),
		CreatedAt:       shop.CreatedAt.Time,
		UpdatedAt:       shop.UpdatedAt.Time,
	}, nil
}

// GetShopsByOwnerID retrieves all shops for a specific owner
func (r *PostgresShopRepository) GetShopsByOwnerID(ctx context.Context, ownerID string) ([]*domain.Shop, error) {
	ownerIDUUID := converter.StringToPgUUID(ownerID)

	shops, err := r.queries.GetShopsByOwnerID(ctx, ownerIDUUID)
	if err != nil {
		log.Println("Error fetching shops by owner ID:", err)
		return nil, fmt.Errorf("failed to get shops by owner ID: %w", err)
	}

	result := make([]*domain.Shop, len(shops))
	for i, shop := range shops {
		result[i] = &domain.Shop{
			ID:              converter.PgUUIDToUUID(shop.ID),
			OwnerID:         converter.PgUUIDToUUID(shop.OwnerID),
			ShopName:        shop.ShopName,
			AvatarURL:       shop.AvatarUrl,
			BannerURL:       shop.BannerUrl,
			ShopDescription: converter.PgTextToStringPtr(shop.ShopDescription),
			AddressID:       converter.PgUUIDToUUID(shop.AddressID),
			Phone:           shop.Phone,
			Email:           shop.Email,
			Rating:          *converter.PgNumericToFloat64Ptr(shop.Rating),
			ActiveAt:        converter.PgTimeToTimePtr(shop.ActiveAt),
			BannedAt:        converter.PgTimeToTimePtr(shop.BannedAt),
			CreatedAt:       shop.CreatedAt.Time,
			UpdatedAt:       shop.UpdatedAt.Time,
		}
	}

	return result, nil
}

// Update updates an existing shop
func (r *PostgresShopRepository) Update(ctx context.Context, shop *domain.Shop) error {
	var rating pgtype.Numeric
	if err := rating.Scan(shop.Rating); err != nil {
		log.Println("Error converting rating:", err)
		return fmt.Errorf("failed to convert rating: %w", err)
	}

	sqlcParams := sqlc.UpdateShopParams{
		ID:              converter.UUIDToPgUUID(shop.ID),
		ShopName:        shop.ShopName,
		AvatarUrl:       shop.AvatarURL,
		BannerUrl:       shop.BannerURL,
		ShopDescription: converter.StringToPgText(shop.ShopDescription),
		AddressID:       converter.UUIDToPgUUID(shop.AddressID),
		Phone:           shop.Phone,
		Email:           shop.Email,
		Rating:          rating,
		ActiveAt:        converter.TimePtrToPgTime(shop.ActiveAt),
		BannedAt:        converter.TimePtrToPgTime(shop.BannedAt),
	}

	updatedShop, err := r.queries.UpdateShop(ctx, sqlcParams)
	if err != nil {
		log.Println("Error updating shop:", err)
		return fmt.Errorf("failed to update shop: %w", err)
	}

	// Update the shop with the updated timestamp
	if updatedShop.UpdatedAt.Valid {
		shop.UpdatedAt = updatedShop.UpdatedAt.Time
	}

	return nil
}

// Delete deletes a shop by ID
func (r *PostgresShopRepository) Delete(ctx context.Context, shopID string) error {
	shopIDUUID := converter.StringToPgUUID(shopID)

	err := r.queries.DeleteShop(ctx, shopIDUUID)
	if err != nil {
		log.Println("Error deleting shop:", err)
		return fmt.Errorf("failed to delete shop: %w", err)
	}

	return nil
}
