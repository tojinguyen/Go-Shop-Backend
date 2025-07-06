package getshop

import (
	"context"
	"time"

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

type GetShopResponse struct {
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
	CreatedAt       string  `json:"created_at"`
}

func (h *Handler) Handle(ctx context.Context, shopID string) (*GetShopResponse, error) {
	shop, err := h.shopRepo.GetShopByID(ctx, shopID)
	if err != nil {
		return nil, err
	}

	response := &GetShopResponse{
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
		CreatedAt:       shop.CreatedAt.Format(time.RFC3339),
	}

	return response, nil
}
