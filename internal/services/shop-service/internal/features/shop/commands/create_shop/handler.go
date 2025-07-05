package createshop

import (
	"context"

	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/service"
)

type Handler struct {
	shopService service.ShopService
}

func NewHandler(shopService service.ShopService) *Handler {
	return &Handler{
		shopService: shopService,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd CreateShopCommand) (*dto.CreateShopResponse, error) {
	// Convert command to DTO
	createCmd := &dto.CreateShopCommand{
		OwnerID:         cmd.OwnerID,
		ShopName:        cmd.ShopName,
		AvatarURL:       cmd.AvatarURL,
		BannerURL:       cmd.BannerURL,
		ShopDescription: cmd.ShopDescription,
		AddressID:       cmd.AddressID,
		Phone:           cmd.Phone,
		Email:           cmd.Email,
	}

	return h.shopService.CreateShop(ctx, createCmd)
}
