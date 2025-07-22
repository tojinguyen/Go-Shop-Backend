package usecase

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/grpc/adapter"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/repository"
)

type OrderUsecase interface {
	CreateOrder(ctx *gin.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error)
}

type orderUsecase struct {
	orderRepo             repository.OrderRepository
	shopServiceAdapter    adapter.ShopServiceAdapter
	productServiceAdapter adapter.ProductServiceAdapter
}

func NewOrderUsecase(orderRepo repository.OrderRepository, shopServiceAdapter adapter.ShopServiceAdapter, productServiceAdapter adapter.ProductServiceAdapter) OrderUsecase {
	return &orderUsecase{orderRepo: orderRepo, shopServiceAdapter: shopServiceAdapter, productServiceAdapter: productServiceAdapter}
}

func (u *orderUsecase) CreateOrder(ctx *gin.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	shopID := req.ShopID

	isShopExists, err := u.shopServiceAdapter.CheckShopExists(ctx, shopID)
	if err != nil {
		return nil, apperror.NewInternal(fmt.Sprintf("Failed to check shop existence: %s", err.Error()))
	}

	if !isShopExists {
		return nil, apperror.NewNotFound("Shop", shopID)
	}

	// Check promotion eligibility
	if *req.PromotionID != "" {
		// Logic to check promotion code validity
		// For now, we assume the promotion code is valid
	}

	//   1. Lấy thông tin UserAddress từ User-Service. Validate.
	//   2. Lấy thông tin Products (giá, tồn kho) từ Product-Service. Validate.
	//   3. Tính ItemsSubtotal từ giá của Product-Service.
	//   4. Nếu có promotionID, gọi Shop-Service để validate và lấy DiscountAmount. Validate.
	//   5. Tính ShippingFee (tạm thời hardcode).
	//   6. Tính TotalAmount.

	//   // (Chưa làm ở bước này, nhưng cần thiết cho hệ thống hoàn chỉnh)
	//   // 7. Gọi Product-Service để TẠM GIỮ kho (Reserve stock). Nếu thất bại -> Lỗi.
	//   // 8. Gọi Payment-Service để ỦY QUYỀN thanh toán (Authorize payment). Nếu thất bại -> Gọi Product-Service để HỦY GIỮ kho -> Lỗi.

	//   9. Bắt đầu DB Transaction:
	//      a. Tạo bản ghi Order với status PENDING.
	//      b. Tạo các bản ghi OrderItem với giá đã chốt.
	//   10. Commit DB Transaction.

	//   // Nếu bước 9 thất bại:
	//   //   a. Rollback Transaction.
	//   //   b. (Nâng cao) Gọi Product-Service để HỦY GIỮ kho.
	//   //   c. (Nâng cao) Gọi Payment-Service để HỦY ủy quyền.
	//   //   d. Trả lỗi Internal Server Error.

	//   11. Publish event `OrderCreatedEvent` vào Message Queue.
	//   12. Trả về `OrderResponse` cho client.

	return nil, nil
}
