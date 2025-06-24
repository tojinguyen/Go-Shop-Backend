package dto

import "time"

// ErrorData represents error information in API responses
type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int `json:"page" form:"page" validate:"min=1"`
	PageSize int `json:"page_size" form:"page_size" validate:"min=1,max=100"`
}

// PaginationResponse represents pagination metadata in responses
type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// SortRequest represents sorting parameters
type SortRequest struct {
	SortBy    string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// FilterRequest represents common filter parameters
type FilterRequest struct {
	Search    string    `json:"search" form:"search"`
	Status    string    `json:"status" form:"status"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`
}

// IDRequest represents a request with an ID parameter
type IDRequest struct {
	ID string `json:"id" uri:"id" validate:"required,uuid"`
}

// DeleteResponse represents a delete operation response
type DeleteResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}
