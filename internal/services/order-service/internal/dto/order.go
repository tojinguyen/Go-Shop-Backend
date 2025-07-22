package dto

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type CreateOrderRequest struct {
	ShopID            string                   `json:"shop_id" binding:"required,uuid"`
	ShippingAddressID string                   `json:"shipping_address_id" binding:"required,uuid"`
	BillingAddressID  string                   `json:"billing_address_id" binding:"required,uuid"`
	PromotionID       *string                  `json:"promotion_id,omitempty" binding:"omitempty,uuid"`
	Note              string                   `json:"note"`
	Items             []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`
}
