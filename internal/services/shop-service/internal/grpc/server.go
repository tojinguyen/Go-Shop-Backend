package grpc

import (
	shopRepo "github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository/shop"
)

type Server struct {
	shopRepo shopRepo.ShopRepository
}

func NewShopGRPCServer(repo shopRepo.ShopRepository) *Server {
	return &Server{
		shopRepo: repo,
	}
}
