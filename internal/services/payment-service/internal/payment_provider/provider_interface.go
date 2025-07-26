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

type CreatePaymentResult struct {
	PayURL string // URL để redirect hoặc hiển thị QR
}

type PaymentProvider interface {
	// GetName trả về tên định danh của provider (vd: "MOMO", "VNPAY")
	GetName() constant.PaymentProviderMethod
	// CreatePayment khởi tạo một giao dịch và trả về thông tin cần thiết cho client
	CreatePayment(ctx context.Context, data PaymentData) (*CreatePaymentResult, error)
	// HandleIPN xử lý callback (Instant Payment Notification) từ cổng thanh toán
	HandleIPN(r *http.Request) (*domain.Payment, error)
}
