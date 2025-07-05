package repository

import (
	"context"
	"fmt"

	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/domain"
)

// PostgresShopRepository implements the ShopRepository interface using PostgreSQL
type PostgresShopRepository struct {
	db *postgresql_infra.PostgreSQLService
}

// NewPostgresShopRepository creates a new PostgreSQL shop repository
func NewPostgresShopRepository(db *postgresql_infra.PostgreSQLService) ShopRepository {
	return &PostgresShopRepository{
		db: db,
	}
}

// Create creates a new shop
func (r *PostgresShopRepository) Create(ctx context.Context, shop *domain.Shop) error {
	query := `
		INSERT INTO shops (
			id, owner_id, shop_name, avatar_url, banner_url, shop_description,
			address_id, phone, email, rating, active_at, banned_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)`

	_, err := r.db.GetPool().Exec(ctx, query,
		shop.ID, shop.OwnerID, shop.ShopName, shop.AvatarURL, shop.BannerURL,
		shop.ShopDescription, shop.AddressID, shop.Phone, shop.Email,
		shop.Rating, shop.ActiveAt, shop.BannedAt, shop.CreatedAt, shop.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create shop: %w", err)
	}

	return nil
}
