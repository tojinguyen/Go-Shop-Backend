package dto

type CreateProductRequest struct {
	ShopID       string  `json:"shop_id" binding:"required,uuid"`
	Name         string  `json:"name" binding:"required,min=5"`
	Description  string  `json:"description" binding:"required"`
	CategoryID   string  `json:"category_id" binding:"required,uuid"`
	Price        float64 `json:"price" binding:"required,gt=0"`
	Currency     string  `json:"currency" binding:"required"`
	ThumbnailURL string  `json:"thumbnail_url" binding:"required,url"`
	Quantity     int     `json:"quantity" binding:"required,gte=0"`
}
