package getshops

import (
	"context"
	"fmt"
	"log"

	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository"
)

// QueryHandler implements GetShopsQueryHandler
type QueryHandler struct {
	shopRepo repository.ShopRepository
}

// NewQueryHandler creates a new GetShopsQueryHandler
func NewQueryHandler(shopRepo repository.ShopRepository) GetShopsQueryHandler {
	return &QueryHandler{
		shopRepo: shopRepo,
	}
}

// Handle processes the GetShopsQuery
func (h *QueryHandler) Handle(ctx context.Context, query GetShopsQuery) ([]*domain.Shop, error) {
	log.Printf("Getting shops for owner ID: %s", query.OwnerID)

	shops, err := h.shopRepo.GetShopsByOwnerID(ctx, query.OwnerID)
	if err != nil {
		log.Printf("Error getting shops by owner ID: %v", err)
		return nil, fmt.Errorf("failed to get shops: %w", err)
	}

	log.Printf("Found %d shops for owner ID: %s", len(shops), query.OwnerID)
	return shops, nil
}
