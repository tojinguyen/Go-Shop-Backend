package usecase

import "github.com/toji-dev/go-shop/internal/services/cart-service/internal/repository"

type cartItemUseCase struct {
	repo repository.CartItemRepository
}

type CartItemUseCase interface {
	AddItemToCart()
	UpdateCartItem()
	RemoveCartItem()
}

func NewCartItemUseCase(repo repository.CartItemRepository) CartItemUseCase {
	return &cartItemUseCase{repo: repo}
}

func (uc *cartItemUseCase) AddItemToCart() {
}

func (uc *cartItemUseCase) UpdateCartItem() {
}

func (uc *cartItemUseCase) RemoveCartItem() {
}
