package usecase

type cartItemUseCase struct {
}

type CartItemUseCase interface {
	AddItemToCart()
	UpdateCartItem()
	RemoveCartItem()
}

func NewCartItemUseCase() CartItemUseCase {
	return &cartItemUseCase{}
}

func (uc *cartItemUseCase) AddItemToCart() {
}

func (uc *cartItemUseCase) UpdateCartItem() {
}

func (uc *cartItemUseCase) RemoveCartItem() {
}