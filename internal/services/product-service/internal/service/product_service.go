package service

import (
	"context"
	"fmt"

	redis_infra "github.com/toji-dev/go-shop/internal/pkg/infra/redis-infra"
	domain "github.com/toji-dev/go-shop/internal/services/product-service/internal/domain/product"
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

func (s *ProductService) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	productRes, err := s.productRepo.Save(ctx, product)

	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return productRes, nil
}
