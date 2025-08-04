package paymentprovider

import (
	"context"
	"net/http"

	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/constant"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/domain"
)

type PaymentData struct {
	OrderID     string
	Amount      int64
	OrderInfo   string
	RedirectURL string
	IPNURL      string
}

type RefundData struct {
	PaymentID             string // ID thanh toán trong hệ thống của bạn
	OrderID               string // ID đơn hàng của bạn
	ProviderTransactionID string // ID giao dịch từ nhà cung cấp (Momo transId)
	Amount                int64
	Reason                string
}

type CreatePaymentResult struct {
	PayURL string // URL để redirect hoặc hiển thị QR
}

type RefundResult struct {
	ProviderRefundID string // ID của giao dịch refund từ nhà cung cấp
	Status           string // Trạng thái refund từ nhà cung cấp
}

type PaymentProvider interface {
	// GetName trả về tên định danh của provider (vd: "MOMO", "VNPAY")
	GetName() constant.PaymentProviderMethod
	// CreatePayment khởi tạo một giao dịch và trả về thông tin cần thiết cho client
	CreatePayment(ctx context.Context, data PaymentData) (*CreatePaymentResult, error)
	// HandleIPN xử lý callback (Instant Payment Notification) từ cổng thanh toán
	HandleIPN(r *http.Request) (*domain.Payment, error)
	// Refund thực hiện hoàn tiền cho giao dịch đã thanh toán
	Refund(ctx context.Context, data RefundData) (*RefundResult, error)
}
