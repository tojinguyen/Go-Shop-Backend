package constant

type PaymentProviderMethod string

const (
	MomoProviderMethod PaymentProviderMethod = "MOMO"
)

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "PENDING"
	PaymentStatusProcessing PaymentStatus = "PROCESSING"
	PaymentStatusSuccess    PaymentStatus = "SUCCESS"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusRefunded   PaymentStatus = "REFUNDED"
)

type PaymentMethod string

const (
	PaymentMethodCOD          PaymentMethod = "COD"
	PaymentMethodCreditCard   PaymentMethod = "CREDIT_CARD"
	PaymentMethodBankTransfer PaymentMethod = "BANK_TRANSFER"
	PaymentMethodEWallet      PaymentMethod = "E_WALLET"
)

type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "PENDING"
	RefundStatusRequested RefundStatus = "REFUND_REQUESTED"
	RefundStatusCompleted RefundStatus = "COMPLETED"
	RefundStatusFailed    RefundStatus = "FAILED"
)
