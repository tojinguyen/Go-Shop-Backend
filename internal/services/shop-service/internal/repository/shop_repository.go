package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/domain"
)

// ShopRepository defines the interface for shop data operations
type ShopRepository interface {
	// Create creates a new shop
	Create(ctx context.Context, shop *domain.Shop) error

	// GetByID retrieves a shop by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Shop, error)

	// GetByOwnerID retrieves shops by owner ID
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.Shop, error)

	// GetByEmail retrieves a shop by email
	GetByEmail(ctx context.Context, email string) (*domain.Shop, error)

	// Update updates an existing shop
	Update(ctx context.Context, shop *domain.Shop) error

	// Delete soft deletes a shop
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves shops with pagination
	List(ctx context.Context, limit, offset int) ([]*domain.Shop, int64, error)

	// Search searches shops by name or description
	Search(ctx context.Context, query string, limit, offset int) ([]*domain.Shop, int64, error)

	// ExistsByEmail checks if a shop with the given email exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByOwnerID checks if the owner already has a shop
	ExistsByOwnerID(ctx context.Context, ownerID uuid.UUID) (bool, error)
}
