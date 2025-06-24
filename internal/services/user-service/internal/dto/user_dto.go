package dto

import "time"

// CreateUserRequest represents the create user request payload
type CreateUserRequest struct {
	Email     string `json:"email" binding:"required" validate:"required,email"`
	Password  string `json:"password" binding:"required" validate:"required,min=8"`
	FirstName string `json:"first_name" binding:"required" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" binding:"required" validate:"required,min=2,max=50"`
	Username  string `json:"username" validate:"omitempty,min=3,max=30"`
	Phone     string `json:"phone" validate:"omitempty,e164"`
	Role      string `json:"role" validate:"omitempty,oneof=user admin shipper"`
}

// UpdateUserRequest represents the update user request payload
type UpdateUserRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,max=50"`
	Username  string `json:"username" validate:"omitempty,min=3,max=30"`
	Phone     string `json:"phone" validate:"omitempty,e164"`
	Avatar    string `json:"avatar" validate:"omitempty,url"`
}

// UpdatePasswordRequest represents the change password request payload
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" validate:"required"`
	NewPassword     string `json:"new_password" binding:"required" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required" validate:"required,eqfield=NewPassword"`
}

// UserResponse represents the user response payload
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Phone     string    `json:"phone,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserListResponse represents the user list response with pagination
type UserListResponse struct {
	Users      []*UserResponse     `json:"users"`
	Pagination *PaginationResponse `json:"pagination"`
}

// UserSearchRequest represents the user search request
type UserSearchRequest struct {
	PaginationRequest
	SortRequest
	Search string `json:"search" form:"search"`
	Role   string `json:"role" form:"role" validate:"omitempty,oneof=user admin shipper"`
	Status string `json:"status" form:"status" validate:"omitempty,oneof=active inactive"`
}

// UserStatusUpdateRequest represents the user status update request
type UserStatusUpdateRequest struct {
	IsActive bool   `json:"is_active" binding:"required"`
	Reason   string `json:"reason" validate:"omitempty,max=255"`
}

// UserProfileUpdateRequest represents the user profile update request
type UserProfileUpdateRequest struct {
	Bio         string `json:"bio" validate:"omitempty,max=500"`
	DateOfBirth string `json:"date_of_birth" validate:"omitempty,datetime=2006-01-02"`
	Gender      string `json:"gender" validate:"omitempty,oneof=male female other"`
	Website     string `json:"website" validate:"omitempty,url"`
}

// UserAddressRequest represents the user address request payload
type UserAddressRequest struct {
	Street     string `json:"street" binding:"required" validate:"required,max=255"`
	City       string `json:"city" binding:"required" validate:"required,max=100"`
	State      string `json:"state" binding:"required" validate:"required,max=100"`
	PostalCode string `json:"postal_code" binding:"required" validate:"required,max=20"`
	Country    string `json:"country" binding:"required" validate:"required,max=100"`
	IsDefault  bool   `json:"is_default"`
}

// UserAddressResponse represents the user address response
type UserAddressResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	PostalCode string    `json:"postal_code"`
	Country    string    `json:"country"`
	IsDefault  bool      `json:"is_default"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
