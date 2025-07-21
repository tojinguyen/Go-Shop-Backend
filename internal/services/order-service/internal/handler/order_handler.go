package handler

import (
	"github.com/gin-gonic/gin"
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

}

func (h *orderHandler) GetOrderByID(c *gin.Context) {

}
