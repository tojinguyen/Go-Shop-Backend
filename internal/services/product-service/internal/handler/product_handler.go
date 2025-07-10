package handler

import (
	"github.com/gin-gonic/gin"
	redis_infra "github.com/toji-dev/go-shop/internal/pkg/infra/redis-infra"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/repository"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(repo repository.ProductRepository, redis *redis_infra.RedisService) *ProductHandler {
	return &ProductHandler{
		productService: service.NewProductService(repo, redis),
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
}
