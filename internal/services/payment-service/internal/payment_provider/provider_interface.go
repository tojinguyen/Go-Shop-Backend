package paymentprovider

import (
	"context"
	"net/http"

	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/constant"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/domain"
)

type PaymentData struct {
	RequestID   string
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

type PaymentStatusResult struct {
	PartnerCode  string `json:"partnerCode"`
	OrderId      string `json:"orderId"`
	RequestId    string `json:"requestId"`
	Amount       int64  `json:"amount"`
	TransId      int64  `json:"transId"`
	PayType      string `json:"payType"`
	ExtraData    string `json:"extraData"`
	Signature    string `json:"signature"`
	ResultCode   int    `json:"resultCode"`
	Message      string `json:"message"`
	ResponseTime string `json:"responseTime"`
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
	// GetPaymentStatus lấy trạng thái của một giao dịch thanh toán
	GetPaymentStatus(ctx context.Context, payment *domain.Payment) (*PaymentStatusResult, error)
}
