package dto

type InitiatePaymentRequest struct {
	OrderID       string  `json:"order_id" binding:"required,uuid"`
	PaymentMethod string  `json:"payment_method" binding:"required,oneof=MOMO VNPAY COD"`
	Amount        float64 `json:"amount" binding:"required,min=1"`
}

type InitiatePaymentResponse struct {
	PaymentID string `json:"payment_id"`
	PayURL    string `json:"pay_url"`
	Message   string `json:"message"`
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
