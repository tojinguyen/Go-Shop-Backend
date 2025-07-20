package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/domain"
	"gorm.io/gorm"
)

type CartItemRepository interface {
	AddItemToCart(ctx *gin.Context, item *domain.CartItem) error
	UpdateCartItem(ctx *gin.Context, item *domain.CartItem) error
	RemoveCartItem(ctx *gin.Context, cartID, productID uuid.UUID) error
}

type cartItemRepository struct {
	db *gorm.DB
}

func NewCartItemRepository(db *gorm.DB) CartItemRepository {
	return &cartItemRepository{db: db}
}

func (r *cartItemRepository) AddItemToCart(ctx *gin.Context, item *domain.CartItem) error {
	return r.db.Create(item).Error
}

func (r *cartItemRepository) UpdateCartItem(ctx *gin.Context, item *domain.CartItem) error {
	return r.db.Model(&domain.CartItem{}).Where("cart_id = ? AND product_id = ?", item.CartID, item.ProductID).Updates(item).Error
}

func (r *cartItemRepository) RemoveCartItem(ctx *gin.Context, cartID, productID uuid.UUID) error {
	return r.db.Delete(&domain.CartItem{}, "cart_id = ? AND product_id = ?", cartID, productID).Error
}
