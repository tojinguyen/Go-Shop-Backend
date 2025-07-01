package domain

type UserProfile struct {
	UserID           string `json:"user_id"`
	Email            string `json:"email"`
	FullName         string `json:"full_name"`
	Birthday         string `json:"birthday"`
	Phone            string `json:"phone"`
	Role             string `json:"role"` // user, shipper
	BannedAt         string `json:"banned_at,omitempty"`
	AvatarURL        string `json:"avatar_url"`
	Gender           string `json:"gender"` //enum: MALE, FEMALE, OTHER
	DefaultAddressID string `json:"default_address_id,omitempty"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}
