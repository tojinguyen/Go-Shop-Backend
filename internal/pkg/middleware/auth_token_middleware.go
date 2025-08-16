package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/pkg/constant"
	jwtService "github.com/toji-dev/go-shop/internal/pkg/jwt"
	"github.com/toji-dev/go-shop/internal/pkg/response"
)

func AuthTokenMiddleware(jwtSvc jwtService.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			appErr := apperror.NewUnauthorized("Authorization header is required")
			response.Unauthorized(c, string(appErr.Code), appErr.Message)
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			appErr := apperror.NewUnauthorized("Invalid authorization header format")
			response.Unauthorized(c, string(appErr.Code), appErr.Message)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			appErr := apperror.NewUnauthorized("Token is missing")
			response.Unauthorized(c, string(appErr.Code), appErr.Message)
			c.Abort()
			return
		}

		// Sử dụng jwtService để xác thực token
		claims, err := jwtSvc.ValidateAccessToken(context.Background(), tokenString)
		if err != nil {
			appErr := apperror.NewUnauthorized("Invalid or expired token")
			response.Unauthorized(c, string(appErr.Code), appErr.Message)
			c.Abort()
			return
		}

		// Nếu token hợp lệ, lưu thông tin user vào context
		c.Set(constant.ContextKeyUserID, claims.UserId)
		c.Set(constant.ContextKeyUserEmail, claims.Email)
		c.Set(constant.ContextKeyUserRole, claims.Role)

		c.Next()
	}
}
