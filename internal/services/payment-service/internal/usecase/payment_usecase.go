package usecase

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/toji-dev/go-shop/internal/pkg/converter"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/constant"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/dto"
	grpc_adapter "github.com/toji-dev/go-shop/internal/services/payment-service/internal/grpc/adapter"
	paymentprovider "github.com/toji-dev/go-shop/internal/services/payment-service/internal/payment_provider"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/repository"
	order_v1 "github.com/toji-dev/go-shop/proto/gen/go/order/v1"
)

type PaymentUseCase interface {
	InitiatePayment(ctx context.Context, userID string, req dto.InitiatePaymentRequest) (*dto.InitiatePaymentResponse, error)
	HandleIPN(ctx context.Context, providerName constant.PaymentProviderMethod, r *http.Request) error
}

type paymentUseCase struct {
	appConfig       *config.AppConfig
	paymentRepo     repository.PaymentRepository
	providerFactory *paymentprovider.PaymentProviderFactory
	orderAdapter    grpc_adapter.OrderServiceAdapter
}

func NewPaymentUsecase(appConfig *config.AppConfig, paymentRepo repository.PaymentRepository, factory *paymentprovider.PaymentProviderFactory, orderAdapter grpc_adapter.OrderServiceAdapter) PaymentUseCase {
	return &paymentUseCase{
		appConfig:       appConfig,
		paymentRepo:     paymentRepo,
		providerFactory: factory,
		orderAdapter:    orderAdapter,
	}
}

func (uc *paymentUseCase) InitiatePayment(ctx context.Context, userID string, req dto.InitiatePaymentRequest) (*dto.InitiatePaymentResponse, error) {
	// 1. Lấy provider từ factory
	paymentProvider, err := uc.providerFactory.GetProvider(constant.PaymentProviderMethod(req.PaymentMethod))
	if err != nil {
		log.Printf("Error getting payment provider %s: %v", req.PaymentMethod, err)
		return nil, err
	}

	orderRequest := &order_v1.GetOrderRequest{
		OrderId: req.OrderID,
	}

	order, err := uc.orderAdapter.GetOrderInfo(ctx, orderRequest)

	if err != nil {
		log.Printf("Error retrieving order info for OrderID %s: %v", req.OrderID, err)
		return nil, fmt.Errorf("could not retrieve order info: %w", err)
	}

	if order == nil {
		log.Printf("Order not found for OrderID %s", req.OrderID)
		return nil, fmt.Errorf("order not found for ID: %s", req.OrderID)
	}

	if order.Order.CustomerId != userID {
		log.Printf("User %s is not authorized to pay for OrderID %s", userID, req.OrderID)
		return nil, fmt.Errorf("user is not authorized to pay for this order")
	}

	amount := float64(order.Order.FinalAmount)
	paymentMethod := strings.ToUpper(req.PaymentMethod)

	// 3. Tạo bản ghi payment trong DB
	params := sqlc.CreatePaymentParams{
		OrderID:         converter.StringToPgUUID(req.OrderID),
		UserID:          converter.StringToPgUUID(userID),
		Amount:          converter.Float64ToPgNumeric(amount),
		Currency:        "VND",
		PaymentMethod:   sqlc.PaymentMethodEWALLET,
		PaymentProvider: converter.StringToPgText(&paymentMethod),
	}

	paymentRecord, err := uc.paymentRepo.CreatePayment(ctx, params)

	if err != nil {
		log.Printf("Error creating payment record for OrderID %s: %v", req.OrderID, err)
		return nil, fmt.Errorf("could not create payment record: %w", err)
	}

	// 4. Gọi provider để tạo link thanh toán
	paymentData := paymentprovider.PaymentData{
		OrderID:     req.OrderID,
		Amount:      int64(amount),
		OrderInfo:   fmt.Sprintf("Thanh_toan_don_hang_%s", req.OrderID),
		IPNURL:      fmt.Sprintf("%s/api/v1/payments/ipn/%s", uc.appConfig.ApiGatewayURL, strings.ToLower(string(paymentProvider.GetName()))),
		RedirectURL: fmt.Sprintf("%s/orders/%s/result", uc.appConfig.FrontendURL, req.OrderID),
	}

	result, err := paymentProvider.CreatePayment(ctx, paymentData)

	if err != nil {
		log.Printf("Error creating payment link for OrderID %s: %v", req.OrderID, err)
		return nil, fmt.Errorf("payment provider failed: %w", err)
	}

	return &dto.InitiatePaymentResponse{
		PaymentID: paymentRecord.ID,
		PayURL:    result.PayURL,
		Message:   "Payment initiated successfully.",
	}, nil
}

func (uc *paymentUseCase) HandleIPN(ctx context.Context, provider constant.PaymentProviderMethod, r *http.Request) error {
	log.Printf("Handling IPN for provider: %s", provider)
	// 1. Lấy provider từ factory
	paymentProvider, err := uc.providerFactory.GetProvider(provider)
	if err != nil {
		log.Printf("Error getting payment provider %s: %v", provider, err)
		return err
	}

	log.Printf("Using provider: %s for IPN handling", paymentProvider.GetName())
	// 2. Dùng provider để xử lý IPN và xác thực
	paymentUpdate, err := paymentProvider.HandleIPN(r)

	if err != nil {
		log.Printf("Error handling IPN for provider %s: %v", provider, err)
		return fmt.Errorf("failed to handle IPN: %w", err)
	}

	// 3. Lấy thông tin payment gốc từ DB
	originalPayment, err := uc.paymentRepo.GetPaymentByOrderID(ctx, paymentUpdate.OrderID)
	if err != nil {
		return fmt.Errorf("original payment record not found for order %s", paymentUpdate.OrderID)
	}

	// 4. Kiểm tra logic nghiệp vụ
	if originalPayment.Status != constant.PaymentStatusPending {
		log.Printf("Payment for OrderID %s already processed. Status: %s. Ignoring IPN.", originalPayment.OrderID, originalPayment.Status)
		return nil
	}
	if originalPayment.Amount != paymentUpdate.Amount {
		log.Printf("Amount mismatch for OrderID %s. DB: %f, Provider: %f", originalPayment.OrderID, originalPayment.Amount, paymentUpdate.Amount)
		return fmt.Errorf("amount mismatch")
	}

	// 5. Cập nhật trạng thái payment
	updateParams := sqlc.UpdatePaymentStatusParams{
		ID:                    converter.StringToPgUUID(originalPayment.ID),
		PaymentStatus:         sqlc.PaymentStatus(paymentUpdate.Status),
		ProviderTransactionID: converter.StringToPgText(paymentUpdate.ProviderTransactionID),
	}
	_, err = uc.paymentRepo.UpdatePaymentStatus(ctx, updateParams)

	if err != nil {
		return fmt.Errorf("failed to update payment status for order %s: %w", originalPayment.OrderID, err)
	}

	// 6. Gửi sự kiện hoặc gọi gRPC tới Order Service để cập nhật trạng thái đơn hàng
	if paymentUpdate.Status == constant.PaymentStatusSuccess {
		log.Printf("Payment for OrderID %s succeeded. Notifying Order Service...", originalPayment.OrderID)
		orderUpdateReq := &order_v1.UpdateOrderStatusRequest{
			OrderId:   originalPayment.OrderID,
			NewStatus: order_v1.OrderStatus_ORDER_STATUS_PROCESSING,
		}
		_, err = uc.orderAdapter.UpdateOrderStatus(ctx, orderUpdateReq)
	} else {
		log.Printf("Payment for OrderID %s failed. Notifying Order Service...", originalPayment.OrderID)

		orderUpdateReq := &order_v1.UpdateOrderStatusRequest{
			OrderId:   originalPayment.OrderID,
			NewStatus: order_v1.OrderStatus_ORDER_STATUS_PAYMENT_FAILED,
		}
		_, _ = uc.orderAdapter.UpdateOrderStatus(ctx, orderUpdateReq) // Có thể bỏ qua lỗi ở đây

		return fmt.Errorf("payment failed for order %s", originalPayment.OrderID)
	}

	return nil
}
