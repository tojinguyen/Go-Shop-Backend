package router

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/services"
)

// RouterConfig holds the configuration for the router
type RouterConfig struct {
	Config           *config.Config
	ServiceContainer *services.ServiceContainer
}

// SetupRoutes sets up all the routes for the user service
func SetupRoutes(serviceContainer *services.ServiceContainer) *gin.Engine {
	cfg := serviceContainer.GetConfig()

	// Set Gin mode based on environment
	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize router
	router := gin.New()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// CORS middleware with configuration
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Set allowed origins from config
		for _, allowedOrigin := range cfg.CORS.AllowedOrigins {
			if allowedOrigin == "*" {
				c.Header("Access-Control-Allow-Origin", "*")
				break
			} else if allowedOrigin == origin {
				c.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}

		// Set other CORS headers from config
		c.Header("Access-Control-Allow-Methods", strings.Join(cfg.CORS.AllowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(cfg.CORS.AllowedHeaders, ", "))

		// Set credentials header if configured
		if cfg.CORS.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// Set max age for preflight cache
		c.Header("Access-Control-Max-Age", "86400") // 24 hours

		// Add expose headers for client access
		c.Header("Access-Control-Expose-Headers", "Authorization, Content-Length, X-CSRF-Token")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize handler factory
	handlerFactory := services.NewHandlerFactory(serviceContainer)

	// Initialize handlers using factory
	authHandler := handlerFactory.CreateAuthHandler()
	profileHandler := handlerFactory.CreateProfileHandler()

	// Health check endpoint with detailed health information
	router.GET("/health", func(c *gin.Context) {
		healthChecker := services.NewHealthChecker(serviceContainer)
		healthInfo := healthChecker.CheckHealth(c.Request.Context())

		// Return appropriate HTTP status based on health
		if healthInfo["status"] == "healthy" {
			c.JSON(200, healthInfo)
		} else {
			c.JSON(503, healthInfo)
		}
	})

	// API versioning
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
			auth.POST("/change-password", authHandler.ChangePassword)
			auth.POST("/validate-access-token", authHandler.ValidateToken)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
		}

		// Protected routes (authentication required)
		protected := v1.Group("/")
		protected.Use(handlerFactory.CreateAuthMiddleware())
		{
			// User profile routes
			profile := protected.Group("/users/profile")
			{
				profile.GET("", profileHandler.GetProfile)
				profile.PUT("", profileHandler.UpdateProfile)
				profile.GET("/:id", profileHandler.GetProfileByID)
				profile.DELETE("", profileHandler.DeleteProfile)
			}
		}
	}

	return router
}
