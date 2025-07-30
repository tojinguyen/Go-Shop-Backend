package grpc_adapter

import (
	"context"
	"log"

	order_v1 "github.com/toji-dev/go-shop/proto/gen/go/order/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderServiceAdapter interface {
	GetOrderInfo(ctx context.Context, req *order_v1.GetOrderRequest) (*order_v1.GetOrderResponse, error)
	UpdateOrderStatus(ctx context.Context, req *order_v1.UpdateOrderStatusRequest) (*order_v1.UpdateOrderStatusResponse, error)
	Close() error
}

type grpcOrderAdapter struct {
	conn   *grpc.ClientConn
	client order_v1.OrderServiceClient
}

func NewGrpcOrderAdapter(orderServiceAddr string) (OrderServiceAdapter, error) {
	log.Printf("Connecting to order service at %s", orderServiceAddr)
	conn, err := grpc.NewClient(orderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to order service: %v", err)
		return nil, err
	}

	client := order_v1.NewOrderServiceClient(conn)

	log.Printf("Successfully connected to order service at %s", orderServiceAddr)

	return &grpcOrderAdapter{
		conn:   conn,
		client: client,
	}, nil
}

func (a *grpcOrderAdapter) GetOrderInfo(ctx context.Context, req *order_v1.GetOrderRequest) (*order_v1.GetOrderResponse, error) {
	return a.client.GetOrder(ctx, req)
}

func (a *grpcOrderAdapter) UpdateOrderStatus(ctx context.Context, req *order_v1.UpdateOrderStatusRequest) (*order_v1.UpdateOrderStatusResponse, error) {
	return a.client.UpdateOrderStatus(ctx, req)
}

func (a *grpcOrderAdapter) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}
