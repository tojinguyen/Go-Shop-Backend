package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/pkg/converter"
	redis_infra "github.com/toji-dev/go-shop/internal/pkg/infra/redis-infra"
	product "github.com/toji-dev/go-shop/internal/services/product-service/internal/domain/product"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/repository"
)

type ProductService struct {
	productRepo  repository.ProductRepository
	redisService *redis_infra.RedisService
}

func NewProductService(productRepo repository.ProductRepository, redisService *redis_infra.RedisService) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		redisService: redisService,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*product.Product, error) {
	shopID, err := uuid.Parse(req.ShopID)
	if err != nil {
		return nil, fmt.Errorf("invalid shop id format")
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category id format")
	}

	price, err := product.NewPrice(req.Price, req.Currency)
	if err != nil {
		return nil, err
	}

	newProduct, err := product.NewProduct(
		shopID.String(),
		categoryID.String(),
		req.Name,
		req.Description,
		converter.StringToUUID(req.ThumbnailURL),
		price,
		req.Quantity,
	)
	if err != nil {
		return nil, err
	}

	if err := s.productRepo.Save(ctx, newProduct); err != nil {
		return nil, fmt.Errorf("could not save product: %w", err)
	}

	return newProduct, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id string) (*product.Product, error) {
	// 1. Validate đầu vào (đơn giản, có thể làm ở handler hoặc đây)
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID format: %w", err)
	}

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		// Service không cần biết lỗi cụ thể là "not found" hay không,
		// nó chỉ cần trả lỗi lên cho handler xử lý.
		return nil, fmt.Errorf("failed to get product from repository: %w", err)
	}

	if product == nil {
		// Repository có thể trả về (nil, nil) để báo hiệu không tìm thấy.
		// Service chuyển nó thành một lỗi rõ ràng hơn.
		return nil, fmt.Errorf("product with ID %s not found", id)
	}

	// 3. Trả về đối tượng domain
	return product, nil
}
