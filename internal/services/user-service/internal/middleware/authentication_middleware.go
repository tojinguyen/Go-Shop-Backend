package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	errorConstants "github.com/your-username/go-shop/internal/services/user-service/internal/pkg/errors"
	jwtService "github.com/your-username/go-shop/internal/services/user-service/internal/pkg/jwt"
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
		c.Set("user_id", claims.UserId)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("user_claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware creates an optional JWT authentication middleware
// This middleware will parse the token if present but won't block the request if missing
func OptionalAuthMiddleware(jwtSvc jwtService.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.Next()
			return
		}

		claims, err := jwtSvc.ValidateAccessToken(c.Request.Context(), token)
		if err != nil {
			// Log the error but don't block the request
			c.Next()
			return
		}

		// Set user information in context if token is valid
		c.Set("user_id", claims.UserId)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("user_claims", claims)

		c.Next()
	}
}

// RoleMiddleware creates a role-based authorization middleware
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "User role not found in context",
				"code":    "ROLE_NOT_FOUND",
				"success": false,
			})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Invalid role type in context",
				"code":    "INVALID_ROLE_TYPE",
				"success": false,
			})
			c.Abort()
			return
		}

		// Check if user role is in allowed roles
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Insufficient permissions",
			"code":    "INSUFFICIENT_PERMISSIONS",
			"success": false,
		})
		c.Abort()
	}
}

// GetUserIDFromContext extracts user ID from gin context
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	id, ok := userID.(string)
	return id, ok
}

// GetUserEmailFromContext extracts user email from gin context
func GetUserEmailFromContext(c *gin.Context) (string, bool) {
	userEmail, exists := c.Get("user_email")
	if !exists {
		return "", false
	}

	email, ok := userEmail.(string)
	return email, ok
}

// GetUserRoleFromContext extracts user role from gin context
func GetUserRoleFromContext(c *gin.Context) (string, bool) {
	userRole, exists := c.Get("user_role")
	if !exists {
		return "", false
	}

	role, ok := userRole.(string)
	return role, ok
}

// GetUserClaimsFromContext extracts user claims from gin context
func GetUserClaimsFromContext(c *gin.Context) (*jwtService.CustomJwtClaims, bool) {
	userClaims, exists := c.Get("user_claims")
	if !exists {
		return nil, false
	}

	claims, ok := userClaims.(*jwtService.CustomJwtClaims)
	return claims, ok
}
