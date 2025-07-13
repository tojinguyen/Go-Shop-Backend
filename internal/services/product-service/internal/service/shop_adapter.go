package service

import (
	"context"

	"github.com/google/uuid"
)

// ShopServiceAdapter là interface để ProductService tương tác.
// Điều này giúp dễ dàng mock khi viết unit test cho ProductService.
type ShopServiceAdapter interface {
	IsShopOwner(ctx context.Context, shopID, userID uuid.UUID) (bool, error)
}

// grpcShopAdapter là implementation của adapter sử dụng gRPC.
type grpcShopAdapter struct {
}

func NewGrpcShopAdapter(shopServiceAddr string) (ShopServiceAdapter, error) {
	return &grpcShopAdapter{}, nil
}

func (a *grpcShopAdapter) IsShopOwner(ctx context.Context, shopID, userID uuid.UUID) (bool, error) {
	return false, nil
}
