package dto

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Email           string `json:"email" binding:"required" validate:"required,email"`
	Password        string `json:"password" binding:"required" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required" validate:"required,eqfield=Password"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=8"`
}

// RefreshTokenRequest represents the refresh token request payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" validate:"required"`
}

type RegisterResponse struct {
	UserID string `json:"user_id"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	ID           string `json:"id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
}

// UserInfo represents user information in auth response
type UserInfo struct {
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
	Token           string `json:"token" binding:"required" validate:"required"`
	NewPassword     string `json:"new_password" binding:"required" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required" validate:"required,eqfield=NewPassword"`
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
