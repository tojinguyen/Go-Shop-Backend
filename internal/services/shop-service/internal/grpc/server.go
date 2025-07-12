package grpc

import (
	shop_v1 "github.com/toji-dev/go-shop/proto/gen/go/shop/v1"
)

type Server struct {
	shop_v1.UnimplementedShopServiceServer // Bắt buộc để tương thích về sau
	shopRepo                               shop.ShopRepository
}

func NewShopGRPCServer(repo shop.ShopRepository) *Server {
	return &Server{
		shopRepo: repo,
	}
}

func (s *Server) CheckShopOwner(ctx context.Context, req *shop_v1.CheckShopOwnerRequest) (*shop_v1.CheckShopOwnerResponse, error) {
	if req.GetShopId() == "" || req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "shop_id and user_id are required")
	}

	shop, err := s.shopRepo.GetShopByID(ctx, req.GetShopId())
	if err != nil {
		// Nếu không tìm thấy shop, coi như không phải chủ sở hữu
		return &shop_v1.CheckShopOwnerResponse{IsOwner: false}, nil
	}

	isOwner := shop.OwnerID.String() == req.GetUserId()

	return &shop_v1.CheckShopOwnerResponse{IsOwner: isOwner}, nil
}
