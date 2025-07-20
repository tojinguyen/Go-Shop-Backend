package repository

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CartRepository interface {
	GetCartByOwnerID(ctx *gin.Context, ownerID uuid.UUID) (*domain.Cart, error)
	Save(ctx *gin.Context, cart *domain.Cart) error
	DeleteCart(ctx *gin.Context, cartID uuid.UUID) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetCartByOwnerID(ctx *gin.Context, ownerID uuid.UUID) (*domain.Cart, error) {
	var cart domain.Cart
	if err := r.db.WithContext(ctx).Preload("Items").First(&cart, "owner_id = ?", ownerID).Error; err != nil {
		log.Printf("Error fetching cart for user %s: %v", ownerID, err)
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.NewNotFound("cart", ownerID.String())
		}
		return nil, apperror.NewInternal(string(apperror.CodeDatabaseError))
	}
	return &cart, nil
}

func (r *cartRepository) Save(ctx *gin.Context, cart *domain.Cart) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Lưu Cart chính
	if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(cart).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(cart).Association("Items").Replace(cart.Items); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *cartRepository) DeleteCart(ctx *gin.Context, cartID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Cart{}, cartID).Error
}
