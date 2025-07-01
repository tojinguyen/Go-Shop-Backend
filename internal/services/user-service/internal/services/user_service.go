package services

import (
	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
)

type UserService struct {
	container *container.ServiceContainer
}

func NewUserService(container *container.ServiceContainer) *UserService {
	return &UserService{
		container: container,
	}
}

func (s *UserService) CreateProfile(ctx *gin.Context, req dto.CreateUserRequest) (domain.UserProfile, error) {
	// Implementation for creating a user profile
	return nil, nil
}
