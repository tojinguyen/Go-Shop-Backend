package dto

import (
	"time"
)

type CreateProductRequest struct {
	ShopID       string  `json:"shop_id" binding:"required,uuid"`
	Name         string  `json:"name" binding:"required,min=5"`
	Description  string  `json:"description" binding:"required"`
	CategoryID   string  `json:"category_id" binding:"uuid"`
	Price        float64 `json:"price" binding:"required,gt=0"`
	Currency     string  `json:"currency" binding:"required"`
	ThumbnailURL string  `json:"thumbnail_url" binding:"required,url"`
	Quantity     int     `json:"quantity" binding:"required,gte=0"`
}

type GetProductsByShopQuery struct {
	ShopID string `form:"shop_id" binding:"required,uuid"`
	Page   int    `form:"page" binding:"required,gte=1"`
	Limit  int    `form:"limit" binding:"required,gte=1,lte=100"`
}

type UpdateProductRequest struct {
	ShopID       string  `json:"shop_id" binding:"required,uuid"`
	Name         string  `json:"name" binding:"required,min=5"`
	Description  string  `json:"description" binding:"required"`
	CategoryID   string  `json:"category_id" binding:"required,uuid"`
	Price        float64 `json:"price" binding:"required,gt=0"`
	Currency     string  `json:"currency" binding:"required"`
	ThumbnailURL string  `json:"thumbnail_url" binding:"required,url"`
	Quantity     int     `json:"quantity" binding:"required,gte=0"`
}

type DeleteProductRequest struct {
	ShopID string `json:"shop_id" binding:"required,uuid"`
}

type ProductResponse struct {
	ID           string    `json:"id"`
	ShopID       string    `json:"shop_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CategoryID   string    `json:"category_id"`
	Price        float64   `json:"price"`
	Currency     string    `json:"currency"`
	Quantity     int       `json:"quantity"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PaginatedProductsResponse struct {
	Data []*ProductResponse `json:"data"`
}
