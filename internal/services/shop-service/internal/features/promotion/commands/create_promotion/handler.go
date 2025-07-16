package createpromotion

import (
	"context"
	"fmt"
	"log"

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

func (h *Handler) Handle(ctx context.Context, shopID uuid.UUID, req CreatePromotionRequest) (*domain.Promotion, error) {
	log.Printf("Creating promotion for shop ID: %f", req.DiscountValue)
	promo := &domain.Promotion{
		ID:                uuid.New(),
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

	err := h.promoRepo.Create(ctx, promo)
	if err != nil {
		return nil, fmt.Errorf("failed to create promotion: %w", err)
	}

	return promo, nil
}
