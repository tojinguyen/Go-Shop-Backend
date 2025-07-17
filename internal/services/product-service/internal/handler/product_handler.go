package handler

import (
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	redis_infra "github.com/toji-dev/go-shop/internal/pkg/infra/redis-infra"
	"github.com/toji-dev/go-shop/internal/pkg/response"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/domain/product"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/repository"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(repo repository.ProductRepository, redis *redis_infra.RedisService, shopService service.ShopServiceAdapter) *ProductHandler {
	return &ProductHandler{
		productService: service.NewProductService(repo, redis, shopService),
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest

	log.Printf("Creating product with request: %+v", req)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		response.BadRequest(c, "VALIDATION_ERROR", "Invalid request body", err.Error())
		return
	}

	productResult, err := h.productService.CreateProduct(c, &req)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		response.InternalServerError(c, "INTERNAL_ERROR", "Failed to create product")
		return
	}

	response.Created(c, "Product created successfully", productResult)
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	// 1. Lấy shopID từ URL (giả sử route là /shops/:shopId/products)
	shopIDStr := c.Param("shopId")
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_SHOP_ID", "Invalid shop ID format", err.Error())
		return
	}

	// 2. Lấy tham số phân trang từ query string
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// 3. Tạo query object và gọi service
	query := dto.GetProductsByShopQuery{
		ShopID: shopID,
		Page:   page,
		Limit:  limit,
	}

	paginatedResult, err := h.productService.GetProductsByShop(c.Request.Context(), query)
	if err != nil {
		response.InternalServerError(c, "GET_PRODUCTS_FAILED", err.Error())
		return
	}

	// 4. Chuyển đổi danh sách domain products sang response DTOs
	productDTOs := make([]*dto.ProductResponse, len(paginatedResult.Products))
	for i, p := range paginatedResult.Products {
		dto := toProductResponse(p)
		productDTOs[i] = &dto
	}

	// 5. Tạo metadata cho phân trang
	meta := response.MetaInfo{
		Page:       page,
		PerPage:    limit,
		Total:      paginatedResult.TotalCount,
		TotalPages: int(math.Ceil(float64(paginatedResult.TotalCount) / float64(limit))),
	}

	// 6. Trả về response hoàn chỉnh
	response.SuccessWithMeta(c, "Products retrieved successfully", productDTOs, &meta)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	// 1. Lấy ID từ URL
	productID := c.Param("id")

	// 2. Gọi Application Service
	product, err := h.productService.GetProductByID(c.Request.Context(), productID)
	if err != nil {
		// 3. Mapping lỗi từ service sang HTTP response
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "PRODUCT_NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "invalid product ID format") {
			response.BadRequest(c, "INVALID_ID", err.Error(), "")
			return
		}

		response.InternalServerError(c, "GET_PRODUCT_FAILED", "Failed to retrieve product")
		return
	}

	// 4. Chuyển đổi domain object sang DTO và trả về
	respDTO := toProductResponse(product)
	response.Success(c, "Product retrieved successfully", respDTO)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	// 1. Lấy ID từ URL
	productID := c.Param("id")
	if productID == "" {
		response.BadRequest(c, "INVALID_ID", "Product ID is required", "")
		return
	}

	// 2. Bind và Validate DTO
	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "VALIDATION_ERROR", "Invalid request body", err.Error())
		return
	}

	// 3. Gọi Application Service
	updatedProduct, err := h.productService.UpdateProduct(c.Request.Context(), productID, req)
	if err != nil {
		// 4. Mapping lỗi
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "PRODUCT_NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "invalid") {
			response.BadRequest(c, "VALIDATION_ERROR", err.Error(), "")
			return
		}

		response.InternalServerError(c, "UPDATE_PRODUCT_FAILED", "Failed to update product")
		return
	}

	// 5. Trả về response thành công
	respDTO := toProductResponse(updatedProduct)
	response.Success(c, "Product updated successfully", respDTO)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	// 1. Lấy ID từ URL
	productID := c.Param("id")

	if productID == "" {
		response.BadRequest(c, "INVALID_ID", "Product ID is required", "")
		return
	}

	// 2. Bind và Validate DTO
	var req dto.DeleteProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "VALIDATION_ERROR", "Invalid request body", err.Error())
		return
	}

	// 2. Gọi Application Service
	err := h.productService.DeleteProduct(c.Request.Context(), productID, req)
	if err != nil {
		// 3. Mapping lỗi
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "PRODUCT_NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "already deleted") {
			response.Conflict(c, "PRODUCT_ALREADY_DELETED", err.Error())
			return
		}
		if strings.Contains(err.Error(), "cannot delete a banned product") {
			response.Forbidden(c, "DELETE_BANNED_PRODUCT_FORBIDDEN", err.Error())
			return
		}

		response.InternalServerError(c, "DELETE_PRODUCT_FAILED", "Failed to delete product")
		return
	}

	// 4. Trả về response thành công
	response.Success(c, "Product deleted successfully", nil)
}

func toProductResponse(p *product.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:           p.ID().String(),
		ShopID:       p.ShopID().String(),
		Name:         p.Name(),
		Description:  *p.Description(),
		CategoryID:   p.CategoryID().String(),
		Price:        p.Price().GetAmount(),
		Currency:     p.Price().GetCurrency(),
		Quantity:     p.Quantity(),
		ThumbnailURL: *p.ThumbnailURL(),
		Status:       string(p.Status()),
		CreatedAt:    p.CreatedAt(),
		UpdatedAt:    p.UpdatedAt(),
	}
}
