package adapter

import (
	"context"

	product_v1 "github.com/toji-dev/go-shop/proto/gen/go/product/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductServiceAdapter interface {
	GetProductsInfo(ctx context.Context, req *product_v1.GetProductsInfoRequest) (*product_v1.GetProductsInfoResponse, error)
	Close() error
}

type grpcProductAdapter struct {
	conn   *grpc.ClientConn
	client product_v1.ProductServiceClient
}

func NewGrpcProductAdapter(productServiceAddr string) (ProductServiceAdapter, error) {
	conn, err := grpc.NewClient(productServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := product_v1.NewProductServiceClient(conn)

	return &grpcProductAdapter{
		conn:   conn,
		client: client,
	}, nil
}

func (a *grpcProductAdapter) GetProductsInfo(ctx context.Context, req *product_v1.GetProductsInfoRequest) (*product_v1.GetProductsInfoResponse, error) {
	return a.client.GetProductsInfo(ctx, req)
}

func (a *grpcProductAdapter) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}
