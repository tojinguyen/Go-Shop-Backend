package domain

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Items     []CartItem
	CreatedAt time.Time
	UpdatedAt time.Time
}
