package service

import (
	"context"

	redis_infra "github.com/toji-dev/go-shop/internal/pkg/infra/redis-infra"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/repository"
	pb "github.com/toji-dev/go-shop/proto/product"
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

func (s *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	return nil, nil
}
