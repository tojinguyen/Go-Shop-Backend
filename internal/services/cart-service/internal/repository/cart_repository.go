package repository

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/pkg/converter"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/domain"
)

type CartRepository interface {
	GetCartByOwnerID(ctx *gin.Context, ownerID uuid.UUID) (*domain.Cart, error)
	Save(ctx *gin.Context, cart *domain.Cart) error
	DeleteCart(ctx *gin.Context, cartID uuid.UUID) error
}

type cartRepository struct {
	db      *postgresql_infra.PostgreSQLService
	queries *sqlc.Queries
}

func NewCartRepository(db *postgresql_infra.PostgreSQLService) CartRepository {
	if db == nil {
		return nil
	}

	queries := sqlc.New(db.GetPool())
	return &cartRepository{
		db:      db,
		queries: queries,
	}
}

func (r *cartRepository) GetCartByOwnerID(ctx *gin.Context, ownerID uuid.UUID) (*domain.Cart, error) {
	ownerIDpg := converter.UUIDToPgUUID(ownerID)

	cart, err := r.queries.GetCartByOwnerID(ctx, ownerIDpg)
	if err != nil {
		return nil, apperror.NewInternal(fmt.Sprintf("failed to get cart by owner ID %s: %v", ownerID, err))
	}

	cartItems, err := r.queries.GetItemsByCartID(ctx, cart.ID)
	if err != nil {
		return nil, apperror.NewInternal(fmt.Sprintf("failed to get cart items by cart ID %s: %v", cart.ID, err))
	}

	domainCart, err := toDomain(&cart)
	if err != nil {
		return nil, apperror.NewInternal(fmt.Sprintf("failed to convert cart to domain: %v", err))
	}

	domainItems, err := toDomainItems(cartItems)
	if err != nil {
		return nil, apperror.NewInternal(fmt.Sprintf("failed to convert cart items to domain: %v", err))
	}
	domainCart.Items = domainItems

	return domainCart, nil
}

func (r *cartRepository) Save(ctx *gin.Context, cart *domain.Cart) error {
	return nil
}

func (r *cartRepository) DeleteCart(ctx *gin.Context, cartID uuid.UUID) error {
	return nil
}

func toDomain(cart *sqlc.Cart) (*domain.Cart, error) {
	if cart == nil {
		return nil, nil
	}

	return &domain.Cart{
		ID:        converter.PgUUIDToUUID(cart.ID),
		UserID:    converter.PgUUIDToUUID(cart.OwnerID),
		CreatedAt: *converter.PgTimeToTimePtr(cart.CreatedAt),
		UpdatedAt: *converter.PgTimeToTimePtr(cart.UpdatedAt),
	}, nil
}

func toDomainItems(items []sqlc.CartItem) ([]domain.CartItem, error) {
	if items == nil {
		return nil, nil
	}

	domainItems := make([]domain.CartItem, len(items))
	for i, item := range items {
		domainItems[i] = domain.CartItem{
			ID:        converter.PgUUIDToUUID(item.ID),
			CartID:    converter.PgUUIDToUUID(item.CartID),
			ProductID: converter.PgUUIDToUUID(item.ProductID),
			Quantity:  int(item.Quantity),
			CreatedAt: *converter.PgTimeToTimePtr(item.CreatedAt),
			UpdatedAt: *converter.PgTimeToTimePtr(item.UpdatedAt),
		}
	}

	return domainItems, nil
}
