package paymentprovider

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/constant"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/domain"
)

type momoProvider struct {
	cfg config.MomoConfig
}

func NewMomoProvider(cfg config.MomoConfig) PaymentProvider {
	return &momoProvider{cfg: cfg}
}

func (p *momoProvider) GetName() constant.PaymentProviderMethod {
	return constant.MomoProviderMethod
}

func (p *momoProvider) CreatePayment(ctx context.Context, data PaymentData) (*CreatePaymentResult, error) {
	// ... (Code tạo request MoMo và signature từ momo_service.go)
	// Chúng ta sẽ copy-paste và điều chỉnh một chút
	requestID := uuid.New().String()
	// MoMo yêu cầu orderId phải là duy nhất cho mỗi lần request
	uniqueOrderID := fmt.Sprintf("%s_%s", data.OrderID, requestID)

	req := &MomoCreatePaymentRequest{
		PartnerCode: p.cfg.PartnerCode,
		RequestID:   requestID,
		Amount:      data.Amount,
		OrderID:     uniqueOrderID,
		OrderInfo:   data.OrderInfo,
		RedirectURL: data.RedirectURL,
		IpnURL:      data.IPNURL,
		RequestType: "captureWallet",
		ExtraData:   "",
		Lang:        "vi",
	}

	rawSignature := fmt.Sprintf("accessKey=%s&amount=%d&extraData=%s&ipnUrl=%s&orderId=%s&orderInfo=%s&partnerCode=%s&redirectUrl=%s&requestId=%s&requestType=%s",
		p.cfg.AccessKey, req.Amount, req.ExtraData, req.IpnURL, req.OrderID, req.OrderInfo, req.PartnerCode, req.RedirectURL, req.RequestID, req.RequestType)
	req.Signature = p.generateSignature(rawSignature)

	reqBody, _ := json.Marshal(req)
	resp, err := http.Post(p.cfg.ApiEndpoint, "application/json", bytes.NewBuffer(reqBody))
	// ... (xử lý response như cũ) ...
	if err != nil {
		log.Printf("Error sending request to MoMo: %v", err)
		return nil, fmt.Errorf("failed to send request to MoMo: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading MoMo response body: %v", err)
		return nil, fmt.Errorf("failed to read MoMo response body: %w", err)
	}

	var momoResp MomoCreatePaymentResponse
	if err := json.Unmarshal(body, &momoResp); err != nil {
		log.Printf("Error unmarshalling MoMo response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal MoMo response: %w", err)
	}

	if momoResp.ResultCode != 0 {
		log.Printf("MoMo returned an error: %s (code: %d)", momoResp.Message, momoResp.ResultCode)
		return nil, fmt.Errorf("momo returned an error: %s (code: %d)", momoResp.Message, momoResp.ResultCode)
	}

	return &CreatePaymentResult{
		PayURL: momoResp.PayURL,
	}, nil
}

func (p *momoProvider) HandleIPN(r *http.Request) (*domain.Payment, error) {
	var req MomoIPNRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode IPN request: %w", err)
	}

	// Lấy lại orderID gốc (bỏ phần unique)
	originalOrderID := strings.Split(req.OrderID, "_")[0]

	// Xác thực chữ ký
	rawSignature := fmt.Sprintf("accessKey=%s&amount=%d&extraData=%s&message=%s&orderId=%s&orderInfo=%s&orderType=%s&partnerCode=%s&payType=%s&requestId=%s&responseTime=%d&resultCode=%d&transId=%d",
		p.cfg.AccessKey, req.Amount, req.ExtraData, req.Message, req.OrderID, req.OrderInfo, req.OrderType,
		req.PartnerCode, req.PayType, req.RequestID, req.ResponseTime, req.ResultCode, req.TransID)

	if !p.verifySignature(rawSignature, req.Signature) {
		return nil, fmt.Errorf("invalid IPN signature")
	}

	// Chuyển đổi thành domain.Payment để usecase xử lý
	payment := &domain.Payment{
		OrderID:               originalOrderID,
		Amount:                float64(req.Amount),
		ProviderTransactionID: &[]string{fmt.Sprintf("%d", req.TransID)}[0],
	}

	if req.ResultCode == 0 {
		payment.Status = constant.PaymentStatusSuccess
	} else {
		payment.Status = constant.PaymentStatusFailed
	}

	return payment, nil
}

type MomoCreatePaymentRequest struct {
	PartnerCode string `json:"partnerCode"`
	RequestID   string `json:"requestId"`
	Amount      int64  `json:"amount"`
	OrderID     string `json:"orderId"`
	OrderInfo   string `json:"orderInfo"`
	RedirectURL string `json:"redirectUrl"`
	IpnURL      string `json:"ipnUrl"`
	RequestType string `json:"requestType"`
	ExtraData   string `json:"extraData"`
	Lang        string `json:"lang"`
	Signature   string `json:"signature"`
}

type MomoCreatePaymentResponse struct {
	PartnerCode  string `json:"partnerCode"`
	OrderID      string `json:"orderId"`
	RequestID    string `json:"requestId"`
	Amount       int64  `json:"amount"`
	ResponseTime int64  `json:"responseTime"`
	Message      string `json:"message"`
	ResultCode   int    `json:"resultCode"`
	PayURL       string `json:"payUrl"`
	Deeplink     string `json:"deeplink"`
}
type MomoIPNRequest struct {
	PartnerCode  string `json:"partnerCode"`
	OrderID      string `json:"orderId"`
	RequestID    string `json:"requestId"`
	Amount       int64  `json:"amount"`
	OrderInfo    string `json:"orderInfo"`
	OrderType    string `json:"orderType"`
	TransID      int64  `json:"transId"`
	ResultCode   int    `json:"resultCode"`
	Message      string `json:"message"`
	PayType      string `json:"payType"`
	ResponseTime int64  `json:"responseTime"`
	ExtraData    string `json:"extraData"`
	Signature    string `json:"signature"`
}

func (p *momoProvider) verifySignature(data, signature string) bool {
	expectedSignature := p.generateSignature(data)
	return expectedSignature == signature
}

func (p *momoProvider) generateSignature(data string) string {
	h := hmac.New(sha256.New, []byte(p.cfg.SecretKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
