package createshop

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository"
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
	// Business logic: Check if owner already has a shop
	exists, err := h.shopRepo.ExistsByOwnerID(ctx, cmd.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check owner existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("owner already has a shop")
	}

	// Business logic: Check if email is already taken
	emailExists, err := h.shopRepo.ExistsByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if emailExists {
		return nil, fmt.Errorf("email already in use")
	}

	// Create shop domain object
	now := time.Now()
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
		ActiveAt:        nil, // Shop starts inactive until approved
		BannedAt:        nil,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Save to repository
	if err := h.shopRepo.Create(ctx, shop); err != nil {
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
