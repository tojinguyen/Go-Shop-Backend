package services

import "github.com/toji-dev/go-shop/internal/services/user-service/internal/container"

type AddressService struct {
	container *container.ServiceContainer
}

func NewAddressService(container *container.ServiceContainer) *AddressService {
	return &AddressService{
		container: container,
	}
}
