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

type PaginatedProducts struct {
	Products   []*product.Product
	TotalCount int64
}

func (s *ProductService) GetProductsByShop(ctx context.Context, query dto.GetProductsByShopQuery) (*PaginatedProducts, error) {
	// 1. Validation logic cho phân trang (đặt ở đây hoặc handler đều được)
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 || query.Limit > 100 { // Giới hạn tối đa 100 sản phẩm/trang
		query.Limit = 20
	}

	// 2. Điều phối Repository để lấy dữ liệu
	// Repository sẽ trả về cả danh sách sản phẩm và tổng số lượng
	products, total, err := s.productRepo.GetByShopID(ctx, query.ShopID, query.Limit, (query.Page-1)*query.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by shop: %w", err)
	}

	return &PaginatedProducts{
		Products:   products,
		TotalCount: total,
	}, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id string, req dto.UpdateProductRequest) (*product.Product, error) {
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
