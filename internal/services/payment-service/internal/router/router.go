package router

import (
	"github.com/gin-gonic/gin"
	common_middleware "github.com/toji-dev/go-shop/internal/pkg/middleware"
	dependency_container "github.com/toji-dev/go-shop/internal/services/payment-service/internal/dependency-container"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/middleware"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func Init(router *gin.Engine, dependencyContainer *dependency_container.DependencyContainer) {
	cfg := dependencyContainer.GetConfig()

	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(common_middleware.ErrorHandler())

	p := ginprometheus.NewPrometheus("go")
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.FullPath()
		if url == "" {
			url = "unknown"
		}
		return url
	}
	p.Use(router)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "payment-service",
		})
	})

	paymentHandler := dependencyContainer.GetPaymentHandler()

	v1 := router.Group("/api/v1")
	{
		payments := v1.Group("/payments")
		payments.POST("/ipn/:provider", paymentHandler.HandleIPN)

		payments.Use(middleware.AuthHeaderMiddleware())
		{
			payments.POST("/initiate", paymentHandler.InitiatePayment)
			payments.POST("/refund", paymentHandler.RefundPayment)
		}
	}
}
