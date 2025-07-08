package deletepromotion

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	repository "github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository/promotion"
)

type Handler struct {
	promoRepo repository.PromotionRepository
}

func NewHandler(promoRepo repository.PromotionRepository) *Handler {
	return &Handler{promoRepo: promoRepo}
}

func (h *Handler) Handle(ctx context.Context, promoID, shopID uuid.UUID) error {
	existingPromo, err := h.promoRepo.GetByID(ctx, promoID.String())
	if err != nil {
		return fmt.Errorf("failed to retrieve existing promotion: %w", err)
	}
	if existingPromo == nil {
		return fmt.Errorf("promotion not found")
	}
	if existingPromo.ShopID != shopID {
		return fmt.Errorf("promotion does not belong to this shop")
	}

	return h.promoRepo.Delete(ctx, promoID.String())
}
