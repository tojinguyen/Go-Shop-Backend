package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/pkg/response"
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

	order, err := h.orderUsecase.CreateOrder(c, request)
	if err != nil {
		response.InternalServerError(c, string(apperror.CodeInternal), fmt.Sprintf("Failed to create order: %s", err.Error()))
		return
	}

	response.Success(c, "Order created successfully", order)
}

func (h *orderHandler) GetOrderByID(c *gin.Context) {

}
