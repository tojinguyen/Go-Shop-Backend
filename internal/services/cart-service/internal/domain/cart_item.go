package domain

import (
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CartID    uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Quantity  int       `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
