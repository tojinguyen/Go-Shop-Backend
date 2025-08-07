package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/usecase"
)

type InboxHandler struct {
	inboxEventUsecase usecase.InboxEventUseCase
}

func NewInboxHandler(inboxEventUsecase usecase.InboxEventUseCase) *InboxHandler {
	return &InboxHandler{
		inboxEventUsecase: inboxEventUsecase,
	}
}

// GetInboxStats - Endpoint để monitor inbox statistics
func (h *InboxHandler) GetInboxStats(c *gin.Context) {
	stats, err := h.inboxEventUsecase.GetInboxStats(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get inbox statistics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// ProcessPendingEvents - Endpoint để manually trigger xử lý pending events
func (h *InboxHandler) ProcessPendingEvents(c *gin.Context) {
	err := h.inboxEventUsecase.ProcessPendingInboxEvents(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process pending events",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pending events processed successfully",
	})
}

// RetryFailedEvents - Endpoint để manually trigger retry failed events
func (h *InboxHandler) RetryFailedEvents(c *gin.Context) {
	err := h.inboxEventUsecase.RetryFailedInboxEvents(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retry failed events",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Failed events retry completed",
	})
}

// CleanupOldEvents - Endpoint để manually trigger cleanup
func (h *InboxHandler) CleanupOldEvents(c *gin.Context) {
	err := h.inboxEventUsecase.CleanupOldEvents(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to cleanup old events",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Old events cleaned up successfully",
	})
}
