package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	common_constant "github.com/toji-dev/go-shop/internal/pkg/constant"
	"github.com/toji-dev/go-shop/internal/pkg/response"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/constant"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/usecase"
)

type PaymentHandler interface {
	InitiatePayment(c *gin.Context)
	HandleIPN(c *gin.Context)
}

type paymentHandler struct {
	paymentUseCase usecase.PaymentUseCase
}

func NewPaymentHandler(paymentUseCase usecase.PaymentUseCase) PaymentHandler {
	return &paymentHandler{paymentUseCase: paymentUseCase}
}

func (h *paymentHandler) InitiatePayment(c *gin.Context) {
	var req dto.InitiatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	userID, _ := c.Get(common_constant.ContextKeyUserID)

	resp, err := h.paymentUseCase.InitiatePayment(c.Request.Context(), userID.(string), req)
	if err != nil {
		log.Printf("Error initiating payment: %v", err)
		response.InternalServerError(c, "PAYMENT_INITIATION_FAILED", err.Error())
		return
	}
	response.Success(c, "Payment link created successfully", resp)
}

func (h *paymentHandler) HandleIPN(c *gin.Context) {
	providerName := c.Param("provider")
	if providerName == "" {
		c.Status(400)
		return
	}

	// Chuyển đổi tên provider sang enum
	provider := constant.PaymentProviderMethod(providerName)

	err := h.paymentUseCase.HandleIPN(c.Request.Context(), provider, c.Request)
	if err != nil {
		log.Printf("Error handling IPN for %s: %v", providerName, err)
		c.Status(500)
		return
	}
	c.Status(204)
}
