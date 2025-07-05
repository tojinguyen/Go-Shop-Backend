package dto

import (
	"github.com/google/uuid"
)

// CreateShopRequest represents the request to create a new shop
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
func (r *CreateShopRequest) ToCommand() *CreateShopCommand {
	ownerID, _ := uuid.Parse(r.OwnerID)
	addressID, _ := uuid.Parse(r.AddressID)

	return &CreateShopCommand{
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

// CreateShopCommand represents the command to create a shop
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

// CreateShopResponse represents the response after creating a shop
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

// ShopDetailsResponse represents detailed shop information
type ShopDetailsResponse struct {
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
	ActiveAt        *string `json:"active_at,omitempty"`
	BannedAt        *string `json:"banned_at,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// UpdateShopRequest represents the request to update a shop
type UpdateShopRequest struct {
	ShopName        *string `json:"shop_name,omitempty" binding:"omitempty,min=1,max=255"`
	AvatarURL       *string `json:"avatar_url,omitempty" binding:"omitempty,url,max=500"`
	BannerURL       *string `json:"banner_url,omitempty" binding:"omitempty,url,max=500"`
	ShopDescription *string `json:"shop_description,omitempty" binding:"omitempty,max=2000"`
	Phone           *string `json:"phone,omitempty" binding:"omitempty,min=10,max=20"`
	Email           *string `json:"email,omitempty" binding:"omitempty,email,max=255"`
}
