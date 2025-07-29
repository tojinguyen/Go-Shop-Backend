package grpc_server

import (
	"context"
	"log"

	"github.com/toji-dev/go-shop/internal/services/order-service/internal/repository"
	order_v1 "github.com/toji-dev/go-shop/proto/gen/go/order/v1"
)

type Server struct {
	order_v1.UnimplementedOrderServiceServer
	orderRepo repository.OrderRepository
}

func NewOrderGRPCServer(orderRepo repository.OrderRepository) *Server {
	return &Server{
		orderRepo: orderRepo,
	}
}

func (s *Server) GetOrder(ctx context.Context, in *order_v1.GetOrderRequest) (*order_v1.GetOrderResponse, error) {
	orderId := in.GetOrderId()

	order, err := s.orderRepo.GetOrderByID(ctx, orderId)

	if err != nil {
		log.Printf("Error retrieving order with ID %s: %v", orderId, err)
		return nil, err
	}

	responseOrder := &order_v1.Order{
		Id:             order.ID,
		CustomerId:     order.OwnerID,
		ShopId:         order.ShopID,
		ShippingFee:    float32(order.ShippingFee),
		DiscountAmount: float32(order.DiscountAmount),
		TotalAmount:    float32(order.TotalAmount),
		FinalAmount:    float32(order.FinalPrice),
		OrderStatus:    string(order.Status),
	}

	return &order_v1.GetOrderResponse{
		Order: responseOrder,
	}, nil
}
