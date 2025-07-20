package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/services/cart-service/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CartRepository interface {
	GetCartByOwnerID(ctx context.Context, ownerID uuid.UUID) (*domain.Cart, error)
	// CreateCart và Update/Save được gộp lại
	Save(ctx context.Context, cart *domain.Cart) error
	DeleteCart(ctx context.Context, cartID uuid.UUID) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetCartByOwnerID(ctx context.Context, ownerID uuid.UUID) (*domain.Cart, error) {
	var cart domain.Cart
	// Preload("Items") sẽ tự động load tất cả CartItem liên quan
	if err := r.db.WithContext(ctx).Preload("Items").First(&cart, "owner_id = ?", ownerID).Error; err != nil {
		log.Printf("Error fetching cart for user %s: %v", ownerID, err)
		if err == gorm.ErrRecordNotFound {
			// Đây là trường hợp không tìm thấy, không phải lỗi hệ thống
			return nil, apperror.NewNotFound("cart", ownerID.String())
		}
		return nil, apperror.NewInternal(string(apperror.CodeDatabaseError))
	}
	return &cart, nil
}

// Save xử lý cả việc tạo mới và cập nhật.
// GORM sẽ tự động xử lý các CartItem liên quan (thêm, sửa, xóa).
func (r *cartRepository) Save(ctx context.Context, cart *domain.Cart) error {
	// Sử dụng transaction để đảm bảo tính toàn vẹn
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Lưu Cart chính
	if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(cart).Error; err != nil {
		tx.Rollback()
		return err
	}

	// GORM's Association Mode để xử lý các Items
	// Xóa hết các item cũ và thay bằng list item mới.
	// Đây là cách đơn giản và hiệu quả nhất để đảm bảo trạng thái đồng bộ.
	if err := tx.Model(cart).Association("Items").Replace(cart.Items); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *cartRepository) DeleteCart(ctx context.Context, cartID uuid.UUID) error {
	// GORM's cascaded delete sẽ tự xóa các items liên quan nếu được cấu hình
	return r.db.WithContext(ctx).Delete(&domain.Cart{}, cartID).Error
}
