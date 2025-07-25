package handler

import (
	"fmt"

	"github.com/toji-dev/go-shop/internal/pkg/constant"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/pkg/response"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/usecase"
)

type OrderHandler interface {
	GetOrdersByOwnerID(c *gin.Context)
	CreateOrder(c *gin.Context)
	GetOrderByID(c *gin.Context)
}

type orderHandler struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(orderUsecase usecase.OrderUsecase) OrderHandler {
	return &orderHandler{orderUsecase: orderUsecase}
}

func (h *orderHandler) GetOrdersByOwnerID(c *gin.Context) {
}

func (h *orderHandler) CreateOrder(c *gin.Context) {
	var request dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, string(apperror.CodeBadRequest), "Invalid request body", err.Error())
		return
	}

	userId, exists := c.Get(constant.ContextKeyUserID)
	if !exists {
		response.Unauthorized(c, string(apperror.CodeUnauthorized), "User not authenticated")
		return
	}

	order, err := h.orderUsecase.CreateOrder(c, userId.(string), request)
	if err != nil {
		response.InternalServerError(c, string(apperror.CodeInternal), fmt.Sprintf("Failed to create order: %s", err.Error()))
		return
	}

	orderResponse := toOrderResponse(order)

	response.Success(c, "Order created successfully", orderResponse)
}

func (h *orderHandler) GetOrderByID(c *gin.Context) {

}

func toOrderResponse(order *domain.Order) *dto.OrderResponse {
	if order == nil {
		return nil
	}

	orderItems := make([]dto.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		orderItems[i] = toOrderItemResponse(&item)
	}

	return &dto.OrderResponse{
		ID:                order.ID,
		ShopID:            order.ShopID,
		ShippingAddressID: order.ShippingAddressID,
		PromotionID:       order.PromotionCode,
		Note:              "",
		DiscountAmount:    order.DiscountAmount,
		TotalAmount:       order.TotalAmount,
		FinalAmount:       order.FinalPrice,
		Status:            string(order.Status),
		CreatedAt:         order.CreatedAt,
		UpdatedAt:         order.UpdatedAt,
		Items:             orderItems,
	}
}

func toOrderItemResponse(item *domain.OrderItem) dto.OrderItemResponse {
	return dto.OrderItemResponse{
		ID:        item.ID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		Price:     item.Price,
	}
}
