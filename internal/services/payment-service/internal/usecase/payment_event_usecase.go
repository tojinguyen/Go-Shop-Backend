package usecase

const (
	BATCH_SIZE = 100
)

type PaymentEventUseCase interface {
	HandlePaymentEvent()
}

type paymentEventUseCase struct {
}

func NewPaymentEventUseCase() PaymentEventUseCase {
	return &paymentEventUseCase{}
}

func (p *paymentEventUseCase) HandlePaymentEvent() {
}
