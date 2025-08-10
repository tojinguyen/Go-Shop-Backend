package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/constant"
	"github.com/toji-dev/go-shop/internal/pkg/response"
)

// AuthorizationMiddleware kiểm tra xem người dùng có vai trò cần thiết không.
func AuthorizationMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy vai trò của người dùng từ context (đã được AuthHeaderMiddleware thiết lập)
		userRole, exists := c.Get(constant.ContextKeyUserRole)
		if !exists {
			response.Forbidden(c, "NO_ROLE_FOUND", "User role not found in token.")
			c.Abort()
			return
		}

		userRoleStr, ok := userRole.(string)
		if !ok {
			response.Forbidden(c, "INVALID_ROLE_FORMAT", "User role has an invalid format.")
			c.Abort()
			return
		}

		// Kiểm tra xem vai trò của người dùng có nằm trong danh sách các vai trò được yêu cầu không
		isAllowed := false
		for _, requiredRole := range requiredRoles {
			if userRoleStr == requiredRole {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			response.Forbidden(c, "INSUFFICIENT_PERMISSIONS", "You do not have permission to perform this action.")
			c.Abort()
			return
		}

		// Nếu có quyền, cho phép đi tiếp
		c.Next()
	}
}
