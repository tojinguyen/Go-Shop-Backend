package createshop

import (
	"github.com/google/uuid"
)

// CreateShopRequest represents the HTTP request to create a new shop
type CreateShopRequest struct {
	OwnerID         string  `json:"owner_id" binding:"required,uuid"`
	ShopName        string  `json:"shop_name" binding:"required,min=1,max=255"`
	AvatarURL       string  `json:"avatar_url" binding:"required,url,max=500"`
	BannerURL       string  `json:"banner_url" binding:"required,url,max=500"`
	ShopDescription *string `json:"shop_description,omitempty" binding:"omitempty,max=2000"`
	AddressID       string  `json:"address_id" binding:"required,uuid"`
	Phone           string  `json:"phone" binding:"required,min=10,max=20"`
	Email           string  `json:"email" binding:"required,email,max=255"`
}

// ToCommand converts the request to a command
func (r *CreateShopRequest) ToCommand() CreateShopCommand {
	ownerID, _ := uuid.Parse(r.OwnerID)
	addressID, _ := uuid.Parse(r.AddressID)

	return CreateShopCommand{
		OwnerID:         ownerID,
		ShopName:        r.ShopName,
		AvatarURL:       r.AvatarURL,
		BannerURL:       r.BannerURL,
		ShopDescription: r.ShopDescription,
		AddressID:       addressID,
		Phone:           r.Phone,
		Email:           r.Email,
	}
}
