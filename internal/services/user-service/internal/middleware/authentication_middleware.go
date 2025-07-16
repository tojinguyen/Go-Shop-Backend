package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/constant"
	errorConstants "github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/errors"
	jwtService "github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/jwt"
)

// AuthMiddleware creates a JWT authentication middleware
func AuthMiddleware(jwtSvc jwtService.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authorization header is required",
				"code":    "MISSING_AUTH_HEADER",
				"success": false,
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid authorization header format. Expected: Bearer <token>",
				"code":    "INVALID_AUTH_FORMAT",
				"success": false,
			})
			c.Abort()
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Token is empty",
				"code":    "EMPTY_TOKEN",
				"success": false,
			})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := jwtSvc.ValidateAccessToken(c.Request.Context(), token)
		if err != nil {
			var errorCode string
			var statusCode int = http.StatusUnauthorized

			switch err {
			case errorConstants.ErrTokenExpired:
				errorCode = "TOKEN_EXPIRED"
			case errorConstants.ErrTokenInvalid:
				errorCode = "TOKEN_INVALID"
			default:
				errorCode = "TOKEN_VALIDATION_FAILED"
			}

			c.JSON(statusCode, gin.H{
				"error":   err.Error(),
				"code":    errorCode,
				"success": false,
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set(constant.ContextKeyUserID, claims.UserId)
		c.Set(constant.ContextKeyUserEmail, claims.Email)
		c.Set(constant.ContextKeyUserRole, claims.Role)
		c.Set("user_claims", claims)

		c.Next()
	}
}

// AuthMiddlewareWithBlacklist creates a JWT authentication middleware with blacklist checking
func AuthMiddlewareWithBlacklist(jwtSvc jwtService.JwtService, authService interface {
	IsTokenBlacklisted(token string) bool
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authorization header is required",
				"code":    "MISSING_AUTH_HEADER",
				"success": false,
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid authorization header format. Expected: Bearer <token>",
				"code":    "INVALID_AUTH_FORMAT",
				"success": false,
			})
			c.Abort()
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Token is empty",
				"code":    "EMPTY_TOKEN",
				"success": false,
			})
			c.Abort()
			return
		}

		// Check if token is blacklisted
		if authService.IsTokenBlacklisted(token) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Token has been invalidated",
				"code":    "TOKEN_BLACKLISTED",
				"success": false,
			})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := jwtSvc.ValidateAccessToken(c.Request.Context(), token)
		if err != nil {
			var errorCode string
			var statusCode int = http.StatusUnauthorized

			switch err {
			case errorConstants.ErrTokenExpired:
				errorCode = "TOKEN_EXPIRED"
			case errorConstants.ErrTokenInvalid:
				errorCode = "TOKEN_INVALID"
			default:
				errorCode = "TOKEN_VALIDATION_FAILED"
			}

			c.JSON(statusCode, gin.H{
				"error":   err.Error(),
				"code":    errorCode,
				"success": false,
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set(constant.ContextKeyUserID, claims.UserId)
		c.Set(constant.ContextKeyUserEmail, claims.Email)
		c.Set(constant.ContextKeyUserRole, claims.Role)
		c.Set("user_claims", claims)
		c.Set("token", token) // Store token for potential use in handlers

		c.Next()
	}
}
