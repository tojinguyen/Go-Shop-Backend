package getshops

import (
	"context"

	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/domain"
)

// GetShopsQuery represents the query to get shops by owner
type GetShopsQuery struct {
	OwnerID string `json:"owner_id" validate:"required,uuid"`
}

// GetShopsQueryHandler handles the GetShopsQuery
type GetShopsQueryHandler interface {
	Handle(ctx context.Context, query GetShopsQuery) ([]*domain.Shop, error)
}
