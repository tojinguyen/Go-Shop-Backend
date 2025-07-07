package domain

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID
	Name        string
	Description string
	ParentID    *uuid.UUID // Nullable, for subcategories
	CreateAt    time.Time
	UpdatedAt   time.Time
}
