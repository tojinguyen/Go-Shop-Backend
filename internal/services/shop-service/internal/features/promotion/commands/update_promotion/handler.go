package updatepromotion

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/domain"
	repository "github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository/promotion"
)

type Handler struct {
	promoRepo repository.PromotionRepository
}

func NewHandler(promoRepo repository.PromotionRepository) *Handler {
	return &Handler{promoRepo: promoRepo}
}

func (h *Handler) Handle(ctx context.Context, promoID, shopID uuid.UUID, req UpdatePromotionRequest) error {
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

	promoToUpdate := &domain.Promotion{
		ID:                promoID,
		ShopID:            shopID,
		PromotionName:     req.PromotionName,
		PromotionType:     domain.PromotionType(req.PromotionType),
		DiscountValue:     req.DiscountValue,
		MaxDiscountAmount: req.MaxDiscountAmount,
		MinPurchaseAmount: req.MinPurchaseAmount,
		UsageLimitPerUser: req.UsageLimitPerUser,
		StartTime:         req.StartTime,
		EndTime:           req.EndTime,
		PromotionStatus:   domain.PromotionStatus(req.PromotionStatus),
	}

	return h.promoRepo.Update(ctx, promoToUpdate)
}
