package router

import (
	"github.com/gin-gonic/gin"
	pkgConstant "github.com/toji-dev/go-shop/internal/pkg/constant" // Import package constant chung
	dependency_container "github.com/toji-dev/go-shop/internal/services/product-service/internal/dependency-container"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/handler"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/middleware"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func SetupRoutes(r *gin.Engine, serviceContainer *dependency_container.DependencyContainer) *gin.Engine {
	p := ginprometheus.NewPrometheus("go")
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.FullPath()
		if url == "" {
			url = "unknown"
		}
		return url
	}
	p.Use(r)

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	productHandler := handler.NewProductHandler(
		serviceContainer.GetProductRepository(),
		serviceContainer.GetRedisService(),
		serviceContainer.GetShopServiceAdapter(),
	)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})

	v1 := r.Group("/api/v1")
	{
		products := v1.Group(("/products"))
		products.Use(middleware.AuthHeaderMiddleware())
		{
			products.POST("", middleware.AuthorizationMiddleware(string(pkgConstant.UserRoleSeller)), productHandler.CreateProduct)
			products.GET("", productHandler.GetProducts)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)
			products.GET("/:id", productHandler.GetProductByID)
		}
	}

	return r
}
