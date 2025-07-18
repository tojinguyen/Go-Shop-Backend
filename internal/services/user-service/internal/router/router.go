package router

import (
	"strings"

	"github.com/gin-gonic/gin"
	common_middleware "github.com/toji-dev/go-shop/internal/pkg/middleware"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/handlers"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/middleware"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

// SetupRoutes sets up all the routes for the user service
func SetupRoutes(serviceContainer container.ServiceContainer) *gin.Engine {
	cfg := serviceContainer.GetConfig()

	// Set Gin mode based on environment
	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize router
	router := gin.New()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(common_middleware.ErrorHandler())

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
	handlerFactory := handlers.NewHandlerFactory(serviceContainer)

	// Initialize handlers using factory
	authHandler := handlerFactory.CreateAuthHandler()
	profileHandler := handlerFactory.CreateProfileHandler()
	addressHandler := handlerFactory.CreateAddressHandler()
	shipperHandler := handlerFactory.CreateShipperHandler()

	// Get AuthService for enhanced middleware
	authService := handlerFactory.GetAuthService()

	// Health check endpoint with detailed health information
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})

	// API versioning
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/reset-password", authHandler.ResetPassword)
			auth.POST("/change-password", authHandler.ChangePassword)
			auth.POST("/validate-access-token", authHandler.ValidateToken)
		}

		// Protected routes (authentication required with blacklist checking)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddlewareWithBlacklist(serviceContainer.GetJWT(), authService))
		{
			// Auth routes that require authentication
			authProtected := protected.Group("/auth")
			{
				authProtected.POST("/logout", authHandler.Logout)
			}

			// User profile routes
			profile := protected.Group("/users/profile")
			{
				profile.POST("", profileHandler.CreateProfile)
				profile.GET("", profileHandler.GetProfile)
				profile.PUT("", profileHandler.UpdateProfile)
				profile.GET("/:id", profileHandler.GetProfileByID)
				profile.DELETE("", profileHandler.DeleteProfile)
			}

			addresses := protected.Group("users/addresses")
			{
				addresses.GET("", addressHandler.GetAddresses)
				addresses.GET("/:id", addressHandler.GetAddressByID)
				addresses.POST("", addressHandler.AddAddress)
				addresses.PUT("/:id", addressHandler.UpdateAddress)
				addresses.DELETE("/:id", addressHandler.DeleteAddress)
				addresses.PUT("/:id/default", addressHandler.SetDefaultAddress)
			}

			shippers := protected.Group("users/shippers")
			{
				shippers.POST("/register", shipperHandler.RegisterShipper)
				shippers.GET("/profile", shipperHandler.GetShipperProfile)
				shippers.GET("/:id/profile", shipperHandler.GetShipperProfileByID)
				shippers.PUT("/profile", shipperHandler.UpdateShipperProfile)
				shippers.DELETE("/profile", shipperHandler.DeleteShipperProfile)
			}
		}
	}

	return router
}
