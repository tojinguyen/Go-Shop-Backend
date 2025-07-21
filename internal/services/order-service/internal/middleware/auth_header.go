package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/constant"
)

func AuthHeaderMiddleware() gin.HandlerFunc {
	log.Println("AuthHeaderMiddleware: Starting authentication header middleware")
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		userRole := c.GetHeader("X-User-Role")
		userEmail := c.GetHeader("X-User-Email")

		log.Printf("AuthHeaderMiddleware: userID=%s, userRole=%s, userEmail=%s", userID, userRole, userEmail)

		// Chỉ cần có UserID là đủ để coi là đã xác thực ở mức cơ bản.
		// Các logic kiểm tra role cụ thể sẽ do từng handler/service đảm nhiệm.
		if userID != "" {
			c.Set(constant.ContextKeyUserID, userID)
			c.Set(constant.ContextKeyUserRole, userRole)
			c.Set(constant.ContextKeyUserEmail, userEmail)
		}

		c.Next()
	}
}
