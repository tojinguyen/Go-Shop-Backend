package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/converter"
	redis_infra "github.com/toji-dev/go-shop/internal/pkg/infra/redis-infra"
	"github.com/toji-dev/go-shop/internal/pkg/response"
	domain "github.com/toji-dev/go-shop/internal/services/product-service/internal/domain/product"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/dto"
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
	var req dto.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "VALIDATION_ERROR", "Invalid request body", err.Error())
		return
	}

	price, err := domain.NewPrice(req.Price, req.Currency)

	if err != nil {
		response.BadRequest(c, "VALIDATION_ERROR", "Invalid price or currency", err.Error())
		return
	}

	categoryID := converter.StringToUUID(req.CategoryID)

	product, err := domain.NewProduct(
		req.ShopID,
		req.Name,
		req.ThumbnailURL,
		req.Description,
		categoryID,
		price,
		req.Quantity,
	)

	if err != nil {
		response.BadRequest(c, "VALIDATION_ERROR", "Invalid product data", err.Error())
		return
	}

	productResult, err := h.productService.CreateProduct(c.Request.Context(), product)
	if err != nil {
		response.InternalServerError(c, "INTERNAL_ERROR", "Failed to create product")
		return
	}

	response.Created(c, "Product created successfully", productResult)
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
}
