package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	shop_v1 "github.com/toji-dev/go-shop/proto/gen/go/shop/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ShopServiceAdapter là interface để ProductService tương tác.
// Điều này giúp dễ dàng mock khi viết unit test cho ProductService.
type ShopServiceAdapter interface {
	IsShopOwner(ctx context.Context, shopID, userID uuid.UUID) (bool, error)
}

// grpcShopAdapter là implementation của adapter sử dụng gRPC.
type grpcShopAdapter struct {
	client shop_v1.ShopServiceClient
}

func NewGrpcShopAdapter(shopServiceAddr string) (ShopServiceAdapter, error) {
	conn, err := grpc.NewClient(shopServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("did not connect to shop service: %w", err)
	}

	client := shop_v1.NewShopServiceClient(conn)
	return &grpcShopAdapter{client: client}, nil
}

func (a *grpcShopAdapter) IsShopOwner(ctx context.Context, shopID, userID uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second) // Set timeout cho RPC call
	defer cancel()

	req := &shop_v1.CheckShopOwnerRequest{
		ShopId: shopID.String(),
		UserId: userID.String(),
	}

	res, err := a.client.CheckShopOwner(ctx, req)
	if err != nil {
		return false, fmt.Errorf("gRPC call to CheckShopOwner failed: %w", err)
	}

	return res.GetIsOwner(), nil
}
