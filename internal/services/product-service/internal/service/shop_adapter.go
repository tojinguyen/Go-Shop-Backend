package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	shop_v1 "github.com/toji-dev/go-shop/proto/gen/go/shop/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ShopServiceAdapter là interface để ProductService tương tác.
// Điều này giúp dễ dàng mock khi viết unit test cho ProductService.
type ShopServiceAdapter interface {
	IsShopOwner(ctx context.Context, shopID, userID uuid.UUID) (bool, error)
	Close() error
}

// grpcShopAdapter là implementation của adapter sử dụng gRPC.
type grpcShopAdapter struct {
	conn   *grpc.ClientConn
	client shop_v1.ShopServiceClient
}

func NewGrpcShopAdapter(shopServiceAddr string) (ShopServiceAdapter, error) {
	conn, err := grpc.NewClient(shopServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := shop_v1.NewShopServiceClient(conn)

	return &grpcShopAdapter{
		conn:   conn,
		client: client,
	}, nil
}

func (a *grpcShopAdapter) IsShopOwner(ctx context.Context, shopID, userID uuid.UUID) (bool, error) {
	req := &shop_v1.CheckShopOwnershipRequest{
		ShopId: shopID.String(),
		UserId: userID.String(),
	}
	res, err := a.client.CheckShopOwnership(ctx, req)
	if err != nil {
		log.Printf("Error calling CheckShopOwnership gRPC: %v", err)
		return false, err
	}

	return res.GetIsOwner(), nil
}

func (a *grpcShopAdapter) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}
