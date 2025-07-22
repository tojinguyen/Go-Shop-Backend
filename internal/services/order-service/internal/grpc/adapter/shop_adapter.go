package adapter

import (
	"context"

	shop_v1 "github.com/toji-dev/go-shop/proto/gen/go/shop/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ShopServiceAdapter interface {
	CheckShopExists(ctx context.Context, shopID string) (bool, error)
	Close() error
}

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

func (a *grpcShopAdapter) CheckShopExists(ctx context.Context, shopID string) (bool, error) {
	req := &shop_v1.CheckShopExistsRequest{
		ShopId: shopID,
	}
	res, err := a.client.CheckShopExists(ctx, req)
	if err != nil {
		return false, err
	}

	return res.GetExists(), nil
}

func (a *grpcShopAdapter) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}
