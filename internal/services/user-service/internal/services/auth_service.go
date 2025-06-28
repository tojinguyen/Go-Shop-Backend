package services

import (
	"context"

	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
)

func Register(ctx context.Context, loginRequest dto.RegisterRequest) (dto.RegisterResponse, error) {
	// Check if user already exists
	return dto.RegisterResponse{}, nil
}
