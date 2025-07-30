package grpc_server

import (
	"context"
	"log"
	"strings"

	"github.com/toji-dev/go-shop/internal/services/order-service/internal/db/sqlc"
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
		OrderStatus:    toProtoOrderStatus(string(order.Status)),
	}

	return &order_v1.GetOrderResponse{
		Order: responseOrder,
	}, nil
}

func (s *Server) UpdateOrderStatus(ctx context.Context, in *order_v1.UpdateOrderStatusRequest) (*order_v1.UpdateOrderStatusResponse, error) {
	orderId := in.GetOrderId()
	statusEnum := in.GetNewStatus()
	statusString := fromProtoOrderStatus(statusEnum)

	_, err := s.orderRepo.UpdateOrderStatus(ctx, orderId, sqlc.OrderStatus(statusString))
	if err != nil {
		log.Printf("Error updating order status for ID %s: %v", orderId, err)
		return nil, err
	}

	return &order_v1.UpdateOrderStatusResponse{
		Success: true,
		Message: "Order status updated successfully",
	}, nil
}

func fromProtoOrderStatus(status order_v1.OrderStatus) string {
	return strings.TrimPrefix(status.String(), "ORDER_STATUS_")
}

func toProtoOrderStatus(status string) order_v1.OrderStatus {
	enumName := "ORDER_STATUS_" + strings.ToUpper(status)
	if val, ok := order_v1.OrderStatus_value[enumName]; ok {
		return order_v1.OrderStatus(val)
	}
	return order_v1.OrderStatus_ORDER_STATUS_UNSPECIFIED
}
