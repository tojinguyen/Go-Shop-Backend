package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/pkg/constant"
	redis_infra "github.com/toji-dev/go-shop/internal/pkg/infra/redis-infra"
	product "github.com/toji-dev/go-shop/internal/services/product-service/internal/domain/product"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/repository"
)

type PaginatedProducts struct {
	Products   []*product.Product
	TotalCount int64
}

type ProductService struct {
	productRepo  repository.ProductRepository
	redisService *redis_infra.RedisService
	shopService  ShopServiceAdapter
}

func NewProductService(productRepo repository.ProductRepository, redisService *redis_infra.RedisService, shopAdapter ShopServiceAdapter) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		redisService: redisService,
		shopService:  shopAdapter,
	}
}

func (s *ProductService) CreateProduct(ctx *gin.Context, req *dto.CreateProductRequest) (*product.Product, error) {
	log.Printf("Creating product for shop %s", req.ShopID)
	userIDCtx := ctx.Value(constant.ContextKeyUserID)
	if userIDCtx == nil {
		log.Fatalf("user_id not found in context")
		return nil, errors.New("unauthorized: user_id not found in context")
	}
	userID, err := uuid.Parse(userIDCtx.(string))
	if err != nil {
		log.Fatalf("invalid user_id format")
		return nil, errors.New("unauthorized: invalid user_id format")
	}

	shopID, err := uuid.Parse(req.ShopID)
	if err != nil {
		log.Fatalf("invalid shop id format")
		return nil, fmt.Errorf("invalid shop id format")
	}

	isOwner, err := s.shopService.IsShopOwner(ctx, shopID, userID)

	if err != nil {
		log.Fatalf("cannot verify shop ownership : %v", err)
		return nil, fmt.Errorf("cannot verify shop ownership: %w", err)
	}
	if !isOwner {
		log.Fatalf("forbidden: you are not the owner of this shop")
		return nil, errors.New("forbidden: you are not the owner of this shop")
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
		req.Name,
		req.ThumbnailURL,
		req.Description,
		categoryID,
		price,
		req.Quantity,
	)

	if err != nil {
		log.Fatalf("could not create new product: %v", err)
		return nil, err
	}

	product, err := s.productRepo.Save(ctx, newProduct)
	if err != nil {
		log.Fatalf("could not create new product: %v", err)
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id string) (*product.Product, error) {
	// 1. Validate đầu vào (đơn giản, có thể làm ở handler hoặc đây)
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID format: %w", err)
	}

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error retrieving product from repository: %v", err)
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

func (s *ProductService) GetProductsByShop(ctx context.Context, query dto.GetProductsByShopQuery) (*PaginatedProducts, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 || query.Limit > 100 {
		query.Limit = 20
	}

	// Parse shopID string to UUID
	shopID, err := uuid.Parse(query.ShopID)
	if err != nil {
		return nil, fmt.Errorf("invalid shop ID format: %w", err)
	}

	products, total, err := s.productRepo.GetByShopID(ctx, shopID, query.Limit, (query.Page-1)*query.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by shop: %w", err)
	}

	return &PaginatedProducts{
		Products:   products,
		TotalCount: total,
	}, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id string, req dto.UpdateProductRequest) (*product.Product, error) {
	userIDCtx := ctx.Value(constant.ContextKeyUserID) // Nên định nghĩa một key cụ thể thay vì string
	if userIDCtx == nil {
		return nil, errors.New("unauthorized: user_id not found in context")
	}
	userID, err := uuid.Parse(userIDCtx.(string))
	if err != nil {
		return nil, errors.New("unauthorized: invalid user_id format")
	}

	shopID, err := uuid.Parse(req.ShopID)
	if err != nil {
		return nil, fmt.Errorf("invalid shop id format")
	}

	isOwner, err := s.shopService.IsShopOwner(ctx, shopID, userID)

	if err != nil {
		return nil, fmt.Errorf("cannot verify shop ownership: %w", err)
	}
	if !isOwner {
		return nil, errors.New("forbidden: you are not the owner of this shop")
	}

	// 1. Lấy Aggregate từ Repository
	existingProduct, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve product for update: %w", err)
	}
	if existingProduct == nil {
		return nil, fmt.Errorf("product with ID %s not found", id)
	}

	// 2. Chuyển đổi dữ liệu từ DTO
	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category id format")
	}
	newPrice, err := product.NewPrice(req.Price, req.Currency)
	if err != nil {
		return nil, err
	}

	// 3. Gọi các phương thức nghiệp vụ trên Domain Object để cập nhật
	if err := existingProduct.ChangeName(req.Name); err != nil {
		return nil, err
	}
	existingProduct.UpdateDescription(req.Description)
	if err := existingProduct.UpdateThumbnail(req.ThumbnailURL); err != nil {
		return nil, err
	}
	existingProduct.ChangeCategory(categoryID)
	if err := existingProduct.ChangePrice(newPrice); err != nil {
		return nil, err
	}
	if err := existingProduct.UpdateQuantity(req.Quantity); err != nil {
		return nil, err
	}

	// 4. Lưu lại Aggregate đã thay đổi
	if err := s.productRepo.Update(ctx, existingProduct); err != nil {
		return nil, fmt.Errorf("failed to save updated product: %w", err)
	}

	return existingProduct, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string, req dto.DeleteProductRequest) error {
	userIDCtx := ctx.Value(constant.ContextKeyUserID) // Nên định nghĩa một key cụ thể thay vì string
	if userIDCtx == nil {
		return errors.New("unauthorized: user_id not found in context")
	}
	userID, err := uuid.Parse(userIDCtx.(string))
	if err != nil {
		return errors.New("unauthorized: invalid user_id format")
	}

	shopID, err := uuid.Parse(req.ShopID)
	if err != nil {
		return fmt.Errorf("invalid shop id format")
	}

	isOwner, err := s.shopService.IsShopOwner(ctx, shopID, userID)

	if err != nil {
		return fmt.Errorf("cannot verify shop ownership: %w", err)
	}
	if !isOwner {
		return errors.New("forbidden: you are not the owner of this shop")
	}

	// 1. Lấy Aggregate từ Repository
	existingProduct, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to retrieve product for deletion: %w", err)
	}
	if existingProduct == nil {
		return fmt.Errorf("product with ID %s not found", id)
	}

	// (Tùy chọn) Kiểm tra quyền hạn ở đây
	// Ví dụ: Lấy user_id từ context và kiểm tra xem có phải chủ shop không.
	// if existingProduct.ShopID() != authorizedUserID { ... }

	// 2. Gọi phương thức nghiệp vụ trên Domain Object
	if err := existingProduct.Delete(); err != nil {
		return err
	}

	if err := s.productRepo.Update(ctx, existingProduct); err != nil {
		return fmt.Errorf("failed to save deleted product state: %w", err)
	}

	return nil
}
