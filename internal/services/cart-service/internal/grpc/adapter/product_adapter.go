package grpc

import (
	"context"

	product_v1 "github.com/toji-dev/go-shop/proto/gen/go/product/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ProductServiceAdapter là interface để ProductService tương tác.
// Điều này giúp dễ dàng mock khi viết unit test cho ProductService.
type ProductServiceAdapter interface {
	GetProductInfo(ctx context.Context, productID string) (*product_v1.GetProductInfoResponse, error)
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

func (a *grpcProductAdapter) GetProductInfo(ctx context.Context, productID string) (*product_v1.GetProductInfoResponse, error) {
	req := &product_v1.GetProductInfoRequest{
		ProductId: productID,
	}
	res, err := a.client.GetProductInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *grpcProductAdapter) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}
