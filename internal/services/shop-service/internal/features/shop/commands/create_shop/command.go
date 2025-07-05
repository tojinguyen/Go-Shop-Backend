package createshop

import "github.com/google/uuid"

type CreateShopCommand struct {
	OwnerID         uuid.UUID `json:"owner_id"`
	ShopName        string    `json:"shop_name"`
	AvatarURL       string    `json:"avatar_url"`
	BannerURL       string    `json:"banner_url"`
	ShopDescription *string   `json:"shop_description,omitempty"`
	AddressID       uuid.UUID `json:"address_id"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
}
