package domain

import (
	"time"

	"github.com/google/uuid"
)

// Shop represents a shop entity
type Shop struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	OwnerID         uuid.UUID  `json:"owner_id" db:"owner_id"`
	ShopName        string     `json:"shop_name" db:"shop_name"`
	AvatarURL       string     `json:"avatar_url" db:"avatar_url"`
	BannerURL       string     `json:"banner_url" db:"banner_url"`
	ShopDescription *string    `json:"shop_description,omitempty" db:"shop_description"`
	AddressID       uuid.UUID  `json:"address_id" db:"address_id"`
	Phone           string     `json:"phone" db:"phone"`
	Email           string     `json:"email" db:"email"`
	Rating          float64    `json:"rating" db:"rating"`
	ActiveAt        *time.Time `json:"active_at,omitempty" db:"active_at"`
	BannedAt        *time.Time `json:"banned_at,omitempty" db:"banned_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// IsActive checks if the shop is currently active
func (s *Shop) IsActive() bool {
	return s.ActiveAt != nil && s.BannedAt == nil
}

// IsBanned checks if the shop is currently banned
func (s *Shop) IsBanned() bool {
	return s.BannedAt != nil
}

// Activate sets the shop as active
func (s *Shop) Activate() {
	now := time.Now()
	s.ActiveAt = &now
	s.BannedAt = nil
	s.UpdatedAt = now
}

// Ban sets the shop as banned
func (s *Shop) Ban() {
	now := time.Now()
	s.BannedAt = &now
	s.UpdatedAt = now
}

// Update updates the shop with new information
func (s *Shop) Update(name, description *string, avatarURL, bannerURL, phone, email *string) {
	if name != nil {
		s.ShopName = *name
	}
	if description != nil {
		s.ShopDescription = description
	}
	if avatarURL != nil {
		s.AvatarURL = *avatarURL
	}
	if bannerURL != nil {
		s.BannerURL = *bannerURL
	}
	if phone != nil {
		s.Phone = *phone
	}
	if email != nil {
		s.Email = *email
	}
	s.UpdatedAt = time.Now()
}
