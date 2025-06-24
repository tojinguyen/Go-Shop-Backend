package router

import (
	"github.com/gin-gonic/gin"
	"github.com/your-username/go-shop/internal/services/user-service/internal/config"
	"github.com/your-username/go-shop/internal/services/user-service/internal/handlers"
	"github.com/your-username/go-shop/internal/services/user-service/internal/middleware"
	jwtService "github.com/your-username/go-shop/internal/services/user-service/internal/pkg/kwt"
)

// RouterConfig holds the configuration for the router
type RouterConfig struct {
	Config     *config.Config
	JWTService jwtService.JwtService
}

// SetupRoutes sets up all the routes for the user service
func SetupRoutes(cfg *config.Config) *gin.Engine {
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

	// CORS middleware (you might want to use a proper CORS package)
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize services
	jwtSvc := jwtService.NewJwtService(cfg)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(jwtSvc, cfg)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "user-service",
			"version": cfg.App.Version,
		})
	})

	// API versioning
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/validate", authHandler.ValidateToken)
		}

		// Protected routes (authentication required)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(jwtSvc))
		{
			// User profile routes
			profile := protected.Group("/profile")
			{
				profile.GET("", authHandler.GetProfile)
			}

			// Logout route
			protected.POST("/auth/logout", authHandler.Logout)
		}
	}

	return router
}

// SetupTestRoutes sets up routes for testing
func SetupTestRoutes() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Load test configuration
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load test configuration: " + err.Error())
	}

	return SetupRoutes(cfg)
}
