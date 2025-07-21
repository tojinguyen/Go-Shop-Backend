package repository

import (
	"fmt"
	"log"

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
	pgCartID := converter.UUIDToPgUUID(cart.ID)
	pgOwnerID := converter.UUIDToPgUUID(cart.UserID)

	tx, err := r.db.BeginTransaction(ctx)
	if err != nil {
		return apperror.NewInternal(fmt.Sprintf("failed to begin transaction: %v", err))
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	upsertCartParams := sqlc.UpsertCartParams{
		ID:      pgCartID,
		OwnerID: pgOwnerID,
	}

	_, upsertCartErr := qtx.UpsertCart(ctx, upsertCartParams)
	if upsertCartErr != nil {
		return apperror.NewInternal(fmt.Sprintf("failed to upsert cart: %v", upsertCartErr))
	}

	dbItems, err := qtx.GetItemsByCartID(ctx, pgCartID)
	if err != nil {
		return apperror.NewInternal(fmt.Sprintf("failed to get cart items by cart ID %s: %v", cart.ID, err))
	}

	dbItemsMap := make(map[uuid.UUID]sqlc.CartItem)
	for _, item := range dbItems {
		dbItemsMap[converter.PgUUIDToUUID(item.ProductID)] = item
	}

	domainItemsMap := make(map[uuid.UUID]domain.CartItem)
	for _, item := range cart.Items {
		domainItemsMap[item.ProductID] = item
	}

	// 3. Xóa những item không còn trong giỏ hàng
	for productID, dbItem := range dbItemsMap {
		if _, exists := domainItemsMap[productID]; !exists {
			if err := qtx.DeleteItemFromCart(ctx, dbItem.ID); err != nil {
				return fmt.Errorf("failed to delete cart item %s: %w", dbItem.ID, err)
			}
		}
	}

	// 4. Upsert (Thêm mới hoặc cập nhật) các item trong giỏ hàng
	for _, domainItem := range cart.Items {
		params := sqlc.UpsertItemInCartParams{
			CartID:    converter.UUIDToPgUUID(cart.ID),
			ProductID: converter.UUIDToPgUUID(domainItem.ProductID),
			ShopID:    converter.UUIDToPgUUID(domainItem.ShopID),
			Quantity:  int32(domainItem.Quantity),
		}

		_, err := qtx.UpsertItemInCart(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to upsert cart item for product %s: %w", domainItem.ProductID, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	log.Printf("Cart %s saved successfully", cart.ID)

	return nil
}

func (r *cartRepository) DeleteCart(ctx *gin.Context, cartID uuid.UUID) error {
	pgCartID := converter.UUIDToPgUUID(cartID)
	// Begin transaction
	tx, err := r.db.BeginTransaction(ctx)
	if err != nil {
		return apperror.NewInternal(fmt.Sprintf("failed to begin transaction: %v", err))
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)
	// Delete cart items
	if err := qtx.DeleteAllItemsFromCart(ctx, pgCartID); err != nil {
		return apperror.NewInternal(fmt.Sprintf("failed to delete cart items for cart ID %s: %v", cartID, err))
	}

	// Delete cart
	if err := qtx.DeleteCart(ctx, pgCartID); err != nil {
		return apperror.NewInternal(fmt.Sprintf("failed to delete cart with ID %s: %v", cartID, err))
	}
	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Cart %s deleted successfully", cartID)
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
