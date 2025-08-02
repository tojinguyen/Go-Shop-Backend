package adapter

import (
	"context"
	"log"

	product_v1 "github.com/toji-dev/go-shop/proto/gen/go/product/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductServiceAdapter interface {
	GetProductsInfo(ctx context.Context, req *product_v1.GetProductsInfoRequest) (*product_v1.GetProductsInfoResponse, error)
	ReserveProducts(ctx context.Context, req *product_v1.ReserveProductsRequest) (*product_v1.ReserveProductsResponse, error)
	UnreserveProducts(ctx context.Context, req *product_v1.UnreserveProductsRequest) (*product_v1.UnreserveProductsResponse, error)
	GetOrderReservationStatus(ctx context.Context, req *product_v1.GetOrderReservationStatusRequest) (*product_v1.GetOrderReservationStatusResponse, error)
	GetOrdersReservationStatus(ctx context.Context, req *product_v1.GetOrdersReservationStatusRequest) (*product_v1.GetOrdersReservationStatusResponse, error)
	Close() error
}

type grpcProductAdapter struct {
	conn   *grpc.ClientConn
	client product_v1.ProductServiceClient
}

func NewGrpcProductAdapter(productServiceAddr string) (ProductServiceAdapter, error) {
	log.Printf("Connecting to product service at %s", productServiceAddr)
	conn, err := grpc.NewClient(productServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to product service: %v", err)
		return nil, err
	}

	client := product_v1.NewProductServiceClient(conn)

	log.Printf("Successfully connected to product service at %s", productServiceAddr)

	return &grpcProductAdapter{
		conn:   conn,
		client: client,
	}, nil
}

func (a *grpcProductAdapter) GetProductsInfo(ctx context.Context, req *product_v1.GetProductsInfoRequest) (*product_v1.GetProductsInfoResponse, error) {
	return a.client.GetProductsInfo(ctx, req)
}

func (a *grpcProductAdapter) ReserveProducts(ctx context.Context, req *product_v1.ReserveProductsRequest) (*product_v1.ReserveProductsResponse, error) {
	return a.client.ReserveProducts(ctx, req)
}

func (a *grpcProductAdapter) UnreserveProducts(ctx context.Context, req *product_v1.UnreserveProductsRequest) (*product_v1.UnreserveProductsResponse, error) {
	return a.client.UnreserveProducts(ctx, req)
}

func (a *grpcProductAdapter) GetOrderReservationStatus(ctx context.Context, req *product_v1.GetOrderReservationStatusRequest) (*product_v1.GetOrderReservationStatusResponse, error) {
	return a.client.GetOrderReservationStatus(ctx, req)
}

func (a *grpcProductAdapter) GetOrdersReservationStatus(ctx context.Context, req *product_v1.GetOrdersReservationStatusRequest) (*product_v1.GetOrdersReservationStatusResponse, error) {
	return a.client.GetOrdersReservationStatus(ctx, req)
}

func (a *grpcProductAdapter) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}
