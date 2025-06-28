package services

import (
	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
)

func Register(ctx *gin.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
	// Check if user already exists

	return dto.RegisterResponse{}, nil
}
