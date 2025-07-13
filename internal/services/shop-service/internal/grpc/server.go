package grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	shopRepo "github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository/shop"
	shop_v1 "github.com/toji-dev/go-shop/proto/gen/go/proto/shop/v1"
)

type Server struct {
	shop_v1.UnimplementedShopServiceServer
	shopRepo shopRepo.ShopRepository
}

func NewShopGRPCServer(repo shopRepo.ShopRepository) *Server {
	return &Server{
		shopRepo: repo,
	}
}

func (s *Server) CheckShopOwnerShip(ctx context.Context, req *shop_v1.CheckShopOwnershipRequest) (*shop_v1.CheckShopOwnershipResponse, error) {
	log.Printf("Received CheckShopOwnership request for ShopID: %s, UserID: %s", req.GetShopId(), req.GetUserId())

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Printf("Invalid UserID format: %s", req.UserId)
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	shopID, err := uuid.Parse(req.ShopId)
	if err != nil {
		log.Printf("Invalid ShopID format: %s", req.ShopId)
		return nil, fmt.Errorf("invalid shop ID format: %w", err)
	}

	shopInfo, err := s.shopRepo.GetShopByID(ctx, shopID.String())
	if err != nil {
		log.Printf("Error retrieving shop with ID %s: %v", shopID, err)
		return nil, fmt.Errorf("error retrieving shop: %w", err)
	}

	if shopInfo == nil {
		log.Printf("Shop with ID %s not found", shopID)
		return nil, fmt.Errorf("shop not found")
	}

	isOwner := shopInfo.OwnerID == userID
	log.Printf("User %s ownership status for Shop %s: %t", userID, shopID, isOwner)

	return &shop_v1.CheckShopOwnershipResponse{
		IsOwner: isOwner,
	}, nil
}
