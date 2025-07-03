package services

import "github.com/toji-dev/go-shop/internal/services/user-service/internal/container"

type ShipperService struct {
	container *container.ServiceContainer
}

func NewShipperService(container *container.ServiceContainer) *ShipperService {
	return &ShipperService{container: container}
}
