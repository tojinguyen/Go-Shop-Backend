package handlers

import (
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/services"
)

// HandlerFactory creates handlers with all necessary dependencies
type HandlerFactory struct {
	container container.ServiceContainer
}

// NewHandlerFactory creates a new handler factory
func NewHandlerFactory(container container.ServiceContainer) *HandlerFactory {
	return &HandlerFactory{
		container: container,
	}
}

// CreateAuthHandler creates an authentication handler
func (hf *HandlerFactory) CreateAuthHandler() *AuthHandler {
	return NewAuthHandler(
		hf.container,
	)
}

// CreateProfileHandler creates a profile handler
func (hf *HandlerFactory) CreateProfileHandler() *ProfileHandler {
	return NewProfileHandler(
		hf.container,
	)
}

func (hf *HandlerFactory) CreateAddressHandler() *AddressHandler {
	return NewAddressHandler(
		hf.container,
	)
}

func (hf *HandlerFactory) CreateShipperHandler() *ShipperHandler {
	return NewShipperHandler(
		hf.container,
	)
}

// GetAuthService returns the auth service for middleware use
func (hf *HandlerFactory) GetAuthService() *services.AuthService {
	return services.NewAuthService(&hf.container)
}
