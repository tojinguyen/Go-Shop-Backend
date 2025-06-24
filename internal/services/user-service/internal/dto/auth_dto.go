package dto

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=8"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Email           string `json:"email" binding:"required" validate:"required,email"`
	Password        string `json:"password" binding:"required" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required" validate:"required,eqfield=Password"`
	FirstName       string `json:"first_name" binding:"required" validate:"required,min=2,max=50"`
	LastName        string `json:"last_name" binding:"required" validate:"required,min=2,max=50"`
	Username        string `json:"username" validate:"omitempty,min=3,max=30"`
	Phone           string `json:"phone" validate:"omitempty,e164"`
}

// RefreshTokenRequest represents the refresh token request payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" validate:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	User         *UserInfo `json:"user"`
}

// UserInfo represents user information in auth response
type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Role      string `json:"role"`
}

// TokenValidationResponse represents the token validation response
type TokenValidationResponse struct {
	Valid     bool   `json:"valid"`
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
	UserRole  string `json:"user_role"`
	ExpiresAt int64  `json:"expires_at"`
	IssuedAt  int64  `json:"issued_at"`
	Issuer    string `json:"issuer"`
}

// UserProfileResponse represents the user profile response
type UserProfileResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required" validate:"required,email"`
}

type ChangePasswordRequest struct {
	Email           string `json:"email" binding:"required" validate:"required,email"`
	CurrentPassword string `json:"current_password" binding:"required" validate:"required,min=8"`
	NewPassword     string `json:"new_password" binding:"required" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required" validate:"required,eqfield=NewPassword"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required" validate:"required,email"`
	OTP   string `json:"otp" binding:"required" validate:"required"`
}
