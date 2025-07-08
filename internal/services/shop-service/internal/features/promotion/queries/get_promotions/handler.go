package getpromotions

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/promotion/dto"
	repository "github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository/promotion"
)

type Handler struct {
	promoRepo repository.PromotionRepository
}

func NewHandler(promoRepo repository.PromotionRepository) *Handler {
	return &Handler{promoRepo: promoRepo}
}

func (h *Handler) Handle(ctx context.Context, shopID uuid.UUID) ([]dto.PromotionResponse, error) {
	promos, err := h.promoRepo.GetByShopID(ctx, shopID.String())
	if err != nil {
		log.Printf("Error getting promotions for shop %s: %v", shopID, err)
		return nil, fmt.Errorf("failed to retrieve promotions: %w", err)
	}

	responses := make([]dto.PromotionResponse, len(promos))
	for i, promo := range promos {
		responses[i] = mapDomainToResponse(promo)
	}

	return responses, nil
}

func mapDomainToResponse(p *domain.Promotion) dto.PromotionResponse {
	return dto.PromotionResponse{
		ID:                p.ID.String(),
		ShopID:            p.ShopID.String(),
		PromotionName:     p.PromotionName,
		PromotionType:     string(p.PromotionType),
		DiscountValue:     p.DiscountValue,
		MaxDiscountAmount: p.MaxDiscountAmount,
		MinPurchaseAmount: p.MinPurchaseAmount,
		UsageLimitPerUser: p.UsageLimitPerUser,
		StartTime:         p.StartTime,
		EndTime:           p.EndTime,
		PromotionStatus:   string(p.PromotionStatus),
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}
}
