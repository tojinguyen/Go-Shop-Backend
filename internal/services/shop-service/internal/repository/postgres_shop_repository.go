package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

// GetByID retrieves a shop by its ID
func (r *PostgresShopRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Shop, error) {
	query := `
		SELECT id, owner_id, shop_name, avatar_url, banner_url, shop_description,
			   address_id, phone, email, rating, active_at, banned_at, created_at, updated_at
		FROM shops 
		WHERE id = $1`

	row := r.db.GetPool().QueryRow(ctx, query, id)

	shop := &domain.Shop{}
	err := row.Scan(
		&shop.ID, &shop.OwnerID, &shop.ShopName, &shop.AvatarURL, &shop.BannerURL,
		&shop.ShopDescription, &shop.AddressID, &shop.Phone, &shop.Email,
		&shop.Rating, &shop.ActiveAt, &shop.BannedAt, &shop.CreatedAt, &shop.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("shop not found")
		}
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}

	return shop, nil
}

// GetByOwnerID retrieves shops by owner ID
func (r *PostgresShopRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.Shop, error) {
	query := `
		SELECT id, owner_id, shop_name, avatar_url, banner_url, shop_description,
			   address_id, phone, email, rating, active_at, banned_at, created_at, updated_at
		FROM shops 
		WHERE owner_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.GetPool().Query(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shops by owner: %w", err)
	}
	defer rows.Close()

	var shops []*domain.Shop
	for rows.Next() {
		shop := &domain.Shop{}
		err := rows.Scan(
			&shop.ID, &shop.OwnerID, &shop.ShopName, &shop.AvatarURL, &shop.BannerURL,
			&shop.ShopDescription, &shop.AddressID, &shop.Phone, &shop.Email,
			&shop.Rating, &shop.ActiveAt, &shop.BannedAt, &shop.CreatedAt, &shop.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan shop: %w", err)
		}
		shops = append(shops, shop)
	}

	return shops, nil
}

// GetByEmail retrieves a shop by email
func (r *PostgresShopRepository) GetByEmail(ctx context.Context, email string) (*domain.Shop, error) {
	query := `
		SELECT id, owner_id, shop_name, avatar_url, banner_url, shop_description,
			   address_id, phone, email, rating, active_at, banned_at, created_at, updated_at
		FROM shops 
		WHERE email = $1`

	row := r.db.GetPool().QueryRow(ctx, query, email)

	shop := &domain.Shop{}
	err := row.Scan(
		&shop.ID, &shop.OwnerID, &shop.ShopName, &shop.AvatarURL, &shop.BannerURL,
		&shop.ShopDescription, &shop.AddressID, &shop.Phone, &shop.Email,
		&shop.Rating, &shop.ActiveAt, &shop.BannedAt, &shop.CreatedAt, &shop.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("shop not found")
		}
		return nil, fmt.Errorf("failed to get shop by email: %w", err)
	}

	return shop, nil
}

// Update updates an existing shop
func (r *PostgresShopRepository) Update(ctx context.Context, shop *domain.Shop) error {
	query := `
		UPDATE shops 
		SET shop_name = $2, avatar_url = $3, banner_url = $4, shop_description = $5,
			phone = $6, email = $7, rating = $8, active_at = $9, banned_at = $10, updated_at = $11
		WHERE id = $1`

	result, err := r.db.GetPool().Exec(ctx, query,
		shop.ID, shop.ShopName, shop.AvatarURL, shop.BannerURL, shop.ShopDescription,
		shop.Phone, shop.Email, shop.Rating, shop.ActiveAt, shop.BannedAt, shop.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update shop: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("shop not found")
	}

	return nil
}

// Delete soft deletes a shop (we can add a deleted_at column later if needed)
func (r *PostgresShopRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM shops WHERE id = $1`

	result, err := r.db.GetPool().Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete shop: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("shop not found")
	}

	return nil
}

// List retrieves shops with pagination
func (r *PostgresShopRepository) List(ctx context.Context, limit, offset int) ([]*domain.Shop, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM shops`
	var total int64
	err := r.db.GetPool().QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get shops with pagination
	query := `
		SELECT id, owner_id, shop_name, avatar_url, banner_url, shop_description,
			   address_id, phone, email, rating, active_at, banned_at, created_at, updated_at
		FROM shops 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.GetPool().Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list shops: %w", err)
	}
	defer rows.Close()

	var shops []*domain.Shop
	for rows.Next() {
		shop := &domain.Shop{}
		err := rows.Scan(
			&shop.ID, &shop.OwnerID, &shop.ShopName, &shop.AvatarURL, &shop.BannerURL,
			&shop.ShopDescription, &shop.AddressID, &shop.Phone, &shop.Email,
			&shop.Rating, &shop.ActiveAt, &shop.BannedAt, &shop.CreatedAt, &shop.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan shop: %w", err)
		}
		shops = append(shops, shop)
	}

	return shops, total, nil
}

// Search searches shops by name or description
func (r *PostgresShopRepository) Search(ctx context.Context, query string, limit, offset int) ([]*domain.Shop, int64, error) {
	searchTerm := "%" + strings.ToLower(query) + "%"

	// Get total count
	countQuery := `
		SELECT COUNT(*) 
		FROM shops 
		WHERE LOWER(shop_name) LIKE $1 OR LOWER(shop_description) LIKE $1`

	var total int64
	err := r.db.GetPool().QueryRow(ctx, countQuery, searchTerm).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get search count: %w", err)
	}

	// Get shops with search and pagination
	searchQuery := `
		SELECT id, owner_id, shop_name, avatar_url, banner_url, shop_description,
			   address_id, phone, email, rating, active_at, banned_at, created_at, updated_at
		FROM shops 
		WHERE LOWER(shop_name) LIKE $1 OR LOWER(shop_description) LIKE $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.GetPool().Query(ctx, searchQuery, searchTerm, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search shops: %w", err)
	}
	defer rows.Close()

	var shops []*domain.Shop
	for rows.Next() {
		shop := &domain.Shop{}
		err := rows.Scan(
			&shop.ID, &shop.OwnerID, &shop.ShopName, &shop.AvatarURL, &shop.BannerURL,
			&shop.ShopDescription, &shop.AddressID, &shop.Phone, &shop.Email,
			&shop.Rating, &shop.ActiveAt, &shop.BannedAt, &shop.CreatedAt, &shop.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan shop: %w", err)
		}
		shops = append(shops, shop)
	}

	return shops, total, nil
}

// ExistsByEmail checks if a shop with the given email exists
func (r *PostgresShopRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM shops WHERE email = $1)`

	var exists bool
	err := r.db.GetPool().QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

// ExistsByOwnerID checks if the owner already has a shop
func (r *PostgresShopRepository) ExistsByOwnerID(ctx context.Context, ownerID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM shops WHERE owner_id = $1)`

	var exists bool
	err := r.db.GetPool().QueryRow(ctx, query, ownerID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check owner existence: %w", err)
	}

	return exists, nil
}
