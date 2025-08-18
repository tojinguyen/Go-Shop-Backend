package adapter

import (
	"context"
	"log"

	shop_v1 "github.com/toji-dev/go-shop/proto/gen/go/shop/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ShopServiceAdapter interface {
	CheckShopExists(ctx context.Context, shopID string) (bool, error)
	CalculatePromotion(ctx context.Context, req *shop_v1.CalculatePromotionRequest) (*shop_v1.CalculatePromotionResponse, error)
	Close() error
}

type grpcShopAdapter struct {
	conn   *grpc.ClientConn
	client shop_v1.ShopServiceClient
}

func NewGrpcShopAdapter(shopServiceAddr string) (ShopServiceAdapter, error) {
	log.Printf("Connecting to shop service at %s", shopServiceAddr)
	conn, err := grpc.NewClient(
		shopServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		log.Printf("Failed to connect to shop service: %v", err)
		return nil, err
	}

	client := shop_v1.NewShopServiceClient(conn)

	log.Printf("Successfully connected to shop service at %s", shopServiceAddr)

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

func (a *grpcShopAdapter) CalculatePromotion(ctx context.Context, req *shop_v1.CalculatePromotionRequest) (*shop_v1.CalculatePromotionResponse, error) {
	return a.client.CalculatePromotion(ctx, req)
}

func (a *grpcShopAdapter) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}
