package usecase

type cartUseCase struct {
}

type CartUseCase interface {
	GetCart()
	DeleteCart()
	ApplyPromotion()
	RemovePromotion()
}

func NewCartUseCase() CartUseCase {
	return &cartUseCase{}
}

func (uc *cartUseCase) GetCart() {
}

func (uc *cartUseCase) DeleteCart() {
}

func (uc *cartUseCase) ApplyPromotion() {
}

func (uc *cartUseCase) RemovePromotion() {
}
