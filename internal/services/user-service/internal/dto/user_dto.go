package dto

// CreateUserProfileRequest represents the create user request payload
type CreateUserProfileRequest struct {
	Email     string `json:"email" binding:"required" validate:"required,email"`
	Password  string `json:"password" binding:"required" validate:"required,min=8"`
	FullName  string `json:"full_name" binding:"required" validate:"required,min=2,max=100"`
	Phone     string `json:"phone" validate:"omitempty,e164"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
	Birthday  string `json:"birthday" validate:"omitempty,datetime=2006-01-02"`
	Gender    string `json:"gender" validate:"omitempty,oneof=male female other"`
}

// UpdateUserProfileRequest represents the update user request payload
type UpdateUserProfileRequest struct {
	FullName  string `json:"full_name" binding:"required" validate:"required,min=2,max=100"`
	Phone     string `json:"phone" validate:"omitempty,e164"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
	Birthday  string `json:"birthday" validate:"omitempty,datetime=2006-01-02"`
}

// UserResponse represents the user response payload
type UserResponse struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	FullName         string `json:"full_name"`
	Birthday         string `json:"birthday,omitempty"`
	Phone            string `json:"phone,omitempty"`
	Avatar           string `json:"avatar,omitempty"`
	Role             string `json:"role"`
	Gender           string `json:"gender"`
	DefaultAddressID string `json:"default_address_id,omitempty"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// PublicUserResponse represents public user profile (limited information)
type PublicUserResponse struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	Avatar    string `json:"avatar,omitempty"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}
