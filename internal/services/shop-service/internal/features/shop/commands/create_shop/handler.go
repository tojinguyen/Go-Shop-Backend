package createshop

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	time_utils "github.com/toji-dev/go-shop/internal/pkg/time"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/domain"
	repository "github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository/shop"
)

type Handler struct {
	shopRepo repository.ShopRepository
}

func NewHandler(shopRepo repository.ShopRepository) *Handler {
	return &Handler{
		shopRepo: shopRepo,
	}
}

type CreateShopResponse struct {
	ID              string  `json:"id"`
	OwnerID         string  `json:"owner_id"`
	ShopName        string  `json:"shop_name"`
	AvatarURL       string  `json:"avatar_url"`
	BannerURL       string  `json:"banner_url"`
	ShopDescription *string `json:"shop_description,omitempty"`
	AddressID       string  `json:"address_id"`
	Phone           string  `json:"phone"`
	Email           string  `json:"email"`
	Rating          float64 `json:"rating"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
}

func (h *Handler) Handle(ctx context.Context, cmd CreateShopCommand) (*CreateShopResponse, error) {
	now := time_utils.GetUtcTime()
	shop := &domain.Shop{
		ID:              uuid.New(),
		OwnerID:         cmd.OwnerID,
		ShopName:        cmd.ShopName,
		AvatarURL:       cmd.AvatarURL,
		BannerURL:       cmd.BannerURL,
		ShopDescription: cmd.ShopDescription,
		AddressID:       cmd.AddressID,
		Phone:           cmd.Phone,
		Email:           cmd.Email,
		Rating:          0.0,
		ActiveAt:        nil,
		BannedAt:        nil,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Save to repository
	if err := h.shopRepo.Create(ctx, shop); err != nil {
		log.Printf("Error creating shop: %v", err)
		return nil, fmt.Errorf("failed to create shop: %w", err)
	}

	// Return response
	return &CreateShopResponse{
		ID:              shop.ID.String(),
		OwnerID:         shop.OwnerID.String(),
		ShopName:        shop.ShopName,
		AvatarURL:       shop.AvatarURL,
		BannerURL:       shop.BannerURL,
		ShopDescription: shop.ShopDescription,
		AddressID:       shop.AddressID.String(),
		Phone:           shop.Phone,
		Email:           shop.Email,
		Rating:          shop.Rating,
		Status:          getShopStatus(shop),
		CreatedAt:       shop.CreatedAt.Format(time.RFC3339),
	}, nil
}

// getShopStatus returns the current status of the shop
func getShopStatus(shop *domain.Shop) string {
	if shop.IsBanned() {
		return "banned"
	}
	if shop.IsActive() {
		return "active"
	}
	return "inactive"
}
