package usecase

type cartUseCase struct {
}

type CartUseCase interface {
}

func NewCartUseCase() CartUseCase {
	return &cartUseCase{}
}
