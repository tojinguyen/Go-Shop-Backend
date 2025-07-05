package updateshop

import (
	"context"

	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/domain"
)

// UpdateShopCommand represents the command to update a shop
type UpdateShopCommand struct {
	ID              string  `json:"id" validate:"required,uuid"`
	ShopName        string  `json:"shop_name" validate:"required,min=1,max=100"`
	AvatarURL       string  `json:"avatar_url" validate:"omitempty,url"`
	BannerURL       string  `json:"banner_url" validate:"omitempty,url"`
	ShopDescription *string `json:"shop_description,omitempty" validate:"omitempty,max=500"`
	AddressID       string  `json:"address_id" validate:"required,uuid"`
	Phone           string  `json:"phone" validate:"required,min=10,max=15"`
	Email           string  `json:"email" validate:"required,email"`
}

// UpdateShopCommandHandler handles the UpdateShopCommand
type UpdateShopCommandHandler interface {
	Handle(ctx context.Context, command UpdateShopCommand) (*domain.Shop, error)
}
