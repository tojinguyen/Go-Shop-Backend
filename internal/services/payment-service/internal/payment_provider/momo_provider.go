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
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/constant"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/dto"
)

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

type MomoRefundRequest struct {
	PartnerCode string `json:"partnerCode"`
	OrderID     string `json:"orderId"`
	RequestID   string `json:"requestId"`
	Amount      int64  `json:"amount"`
	TransID     int64  `json:"transId"`
	Lang        string `json:"lang"`
	Description string `json:"description"`
	Signature   string `json:"signature"`
}

type MomoRefundResponse struct {
	PartnerCode string `json:"partnerCode"`
	OrderID     string `json:"orderId"`
	RequestID   string `json:"requestId"`
	Amount      int64  `json:"amount"`
	ResultCode  int    `json:"resultCode"`
	Message     string `json:"message"`
	TransID     int64  `json:"transId"`
}

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
	requestID := data.RequestID
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
		p.cfg.AccessKey,
		req.Amount,
		req.ExtraData,
		req.IpnURL,
		req.OrderID,
		req.OrderInfo,
		p.cfg.PartnerCode,
		req.RedirectURL,
		req.RequestID,
		req.RequestType,
	)
	req.Signature = p.generateSignature(rawSignature)

	reqBody, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error marshalling MoMo request: %v", err)
		return nil, fmt.Errorf("failed to marshal MoMo request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, "POST", p.cfg.ApiEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Error creating MoMo HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create MoMo HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpRequest)
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

	log.Printf("[MOMO DEBUG] MoMo Response Body: %s", string(body))

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
	var req dto.MomoIPNRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[MOMO IPN] Error decoding IPN request: %v", err)
		return nil, fmt.Errorf("failed to decode IPN request: %w", err)
	}

	log.Printf("[MOMO IPN] Received IPN: OrderID=%s, TransID=%d, ResultCode=%d", req.OrderID, req.TransID, req.ResultCode)

	// Lấy lại orderID gốc (bỏ phần unique)
	originalOrderID := strings.Split(req.OrderID, "_")[0]

	rawSignature := fmt.Sprintf("accessKey=%s&amount=%d&extraData=%s&message=%s&orderId=%s&orderInfo=%s&orderType=%s&partnerCode=%s&payType=%s&requestId=%s&responseTime=%d&resultCode=%d&transId=%d",
		p.cfg.AccessKey,
		req.Amount,
		req.ExtraData,
		req.Message,
		req.OrderID,
		req.OrderInfo,
		req.OrderType,
		req.PartnerCode,
		req.PayType,
		req.RequestID,
		req.ResponseTime,
		req.ResultCode,
		req.TransID,
	)

	log.Printf("[MOMO IPN] Raw Signature for verification: %s", rawSignature)

	if !p.verifySignature(rawSignature, req.Signature) {
		log.Printf("[MOMO IPN] Signature verification failed. Expected vs Received: %s vs %s",
			p.generateSignature(rawSignature), req.Signature)
		return nil, fmt.Errorf("invalid IPN signature")
	}

	log.Printf("[MOMO IPN] Signature verification successful")

	// Chuyển đổi thành domain.Payment để usecase xử lý
	log.Printf("[MOMO IPN] Request IPN: %+v", req)
	log.Printf("[MOMO IPN] Processing payment for OrderID: %d", req.TransID)
	providerTransID := fmt.Sprintf("%d", req.TransID)
	log.Printf("[MOMO IPN] Provider Transaction ID: %s", providerTransID)
	payment := &domain.Payment{
		OrderID:               originalOrderID,
		Amount:                float64(req.Amount),
		ProviderTransactionID: &providerTransID,
	}

	if req.ResultCode == 0 {
		payment.Status = constant.PaymentStatusSuccess
		log.Printf("[MOMO IPN] Payment successful for OrderID: %s", originalOrderID)
	} else {
		payment.Status = constant.PaymentStatusFailed
		log.Printf("[MOMO IPN] Payment failed for OrderID: %s, Message: %s", originalOrderID, req.Message)
	}

	return payment, nil
}

func (p *momoProvider) Refund(ctx context.Context, data RefundData) (*RefundResult, error) {
	requestID := uuid.New().String()
	transID, err := strconv.ParseInt(data.ProviderTransactionID, 10, 64)
	if err != nil {
		log.Printf("[MOMO REFUND] Invalid provider transaction ID format: %v", err)
		return nil, fmt.Errorf("invalid provider transaction ID format: %w", err)
	}

	req := &MomoRefundRequest{
		PartnerCode: p.cfg.PartnerCode,
		OrderID:     data.OrderID,
		RequestID:   requestID,
		Amount:      data.Amount,
		TransID:     transID,
		Lang:        "vi",
		Description: data.Reason,
	}

	rawSignature := fmt.Sprintf(
		"accessKey=%s&amount=%d&description=%s&orderId=%s&partnerCode=%s&requestId=%s&transId=%d",
		p.cfg.AccessKey,
		req.Amount,
		req.Description,
		req.OrderID,
		p.cfg.PartnerCode,
		req.RequestID,
		req.TransID,
	)

	req.Signature = p.generateSignature(rawSignature)

	reqBody, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error marshalling MoMo refund request: %v", err)
		return nil, fmt.Errorf("failed to marshal MoMo refund request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, "POST", p.cfg.ApiRefundEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Error creating MoMo refund HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create MoMo refund HTTP request: %w", err)
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(httpRequest)
	if err != nil {
		log.Printf("Error sending refund request to MoMo: %v", err)
		return nil, fmt.Errorf("failed to send refund request to MoMo: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading MoMo refund response body: %v", err)
		return nil, fmt.Errorf("failed to read MoMo refund response body: %w", err)
	}

	var momoResp MomoRefundResponse
	if err := json.Unmarshal(body, &momoResp); err != nil {
		log.Printf("Error unmarshalling MoMo refund response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal MoMo refund response: %w", err)
	}

	if momoResp.ResultCode != 0 && momoResp.ResultCode != 10 {
		log.Printf("MoMo refund returned an error: %s (code: %d)", momoResp.Message, momoResp.ResultCode)
		return nil, fmt.Errorf("momo refund failed with code %d: %s", momoResp.ResultCode, momoResp.Message)
	}

	return &RefundResult{
		ProviderRefundID: fmt.Sprintf("%d", momoResp.TransID),
		Status:           string(sqlc.RefundStatusCOMPLETED),
	}, nil
}

func (p *momoProvider) GetPaymentStatus(ctx context.Context, payment *domain.Payment) (*PaymentStatusResult, error) {
	requestID := uuid.New().String()
	orderID := fmt.Sprintf("%s_%s", payment.OrderID, requestID)

	rawData := fmt.Sprintf(
		"accessKey=%s&orderId=%s&partnerCode=%s&requestId=%s",
		p.cfg.AccessKey,
		orderID,
		p.cfg.PartnerCode,
		requestID,
	)

	signature := p.generateSignature(rawData)

	// Tạo payload gửi MoMo
	req := map[string]string{
		"partnerCode": p.cfg.PartnerCode,
		"orderId":     orderID,
		"requestId":   requestID,
		"lang":        "vi",
		"signature":   signature,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	queryEndpoint := p.cfg.ApiGetStatusEndpoint
	httpRequest, err := http.NewRequestWithContext(ctx, "POST", queryEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Error creating MoMo query HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create MoMo query HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpRequest)
	if err != nil {
		log.Printf("Error sending query request to MoMo: %v", err)
		return nil, fmt.Errorf("failed to send query request to MoMo: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading MoMo query response body: %v", err)
		return nil, fmt.Errorf("failed to read MoMo query response body: %w", err)
	}

	log.Printf("[MOMO DEBUG] MoMo Query Response Body: %s", string(body))

	// Parse response
	var momoResp PaymentStatusResult
	if err := json.Unmarshal(body, &momoResp); err != nil {
		log.Printf("Error unmarshalling MoMo query response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal MoMo query response: %w", err)
	}

	// Kiểm tra kết quả
	if momoResp.ResultCode != 0 {
		log.Printf("MoMo query returned an error: %s (code: %d)", momoResp.Message, momoResp.ResultCode)
		return nil, fmt.Errorf("momo query returned an error: %s (code: %d)", momoResp.Message, momoResp.ResultCode)
	}

	momoResp.PaymentStatus = MapResultCodeToStatus(momoResp.ResultCode)

	return &momoResp, nil
}

func (p *momoProvider) verifySignature(data, signature string) bool {
	expectedSignature := p.generateSignature(data)
	isValid := expectedSignature == signature

	if !isValid {
		log.Printf("[MOMO DEBUG] Signature mismatch - Expected: %s, Got: %s", expectedSignature, signature)
	}

	return isValid
}

func (p *momoProvider) generateSignature(data string) string {
	h := hmac.New(sha256.New, []byte(p.cfg.SecretKey))
	h.Write([]byte(data))
	signature := hex.EncodeToString(h.Sum(nil))

	log.Printf("[MOMO DEBUG] Generated signature for data '%s': %s", data, signature)

	return signature
}
