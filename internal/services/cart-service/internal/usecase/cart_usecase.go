package usecase

import "github.com/toji-dev/go-shop/internal/services/cart-service/internal/repository"

type cartUseCase struct {
	repo repository.CartRepository
}

type CartUseCase interface {
	GetCart()
	DeleteCart()
	ApplyPromotion()
	RemovePromotion()
}

func NewCartUseCase(repo repository.CartRepository) CartUseCase {
	return &cartUseCase{repo: repo}
}

func (uc *cartUseCase) GetCart() {
}

func (uc *cartUseCase) DeleteCart() {
}

func (uc *cartUseCase) ApplyPromotion() {
}

func (uc *cartUseCase) RemovePromotion() {
}
