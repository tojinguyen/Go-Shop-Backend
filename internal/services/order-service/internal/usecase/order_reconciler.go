package usecase

import (
	"context"
	"log"
	"time"

	time_utils "github.com/toji-dev/go-shop/internal/pkg/time"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/grpc/adapter"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/repository"
	product_v1 "github.com/toji-dev/go-shop/proto/gen/go/product/v1"
)

const (
	STALE_ORDER_THRESHOLD = 10 * time.Minute
	BATCH_SIZE            = 100
	ORDER_IN_ONE_CALL     = 10
)

type OrderReconciler struct {
	orderRepo      repository.OrderRepository
	productAdapter adapter.ProductServiceAdapter
}

func NewOrderReconciler(orderRepo repository.OrderRepository, productAdapter adapter.ProductServiceAdapter) *OrderReconciler {
	return &OrderReconciler{
		orderRepo:      orderRepo,
		productAdapter: productAdapter,
	}
}

func (r *OrderReconciler) ReconcilePendingOrders() {
	ctx := context.Background()

	log.Println("[OrderReconciler] Starting reconciliation of pending orders...")

	olderThan := time_utils.GetUtcTime().Add(-STALE_ORDER_THRESHOLD)
	orders, err := r.orderRepo.GetStaleOrders(ctx, olderThan, BATCH_SIZE)

	if err != nil {
		log.Printf("[OrderReconciler] Error fetching stale orders: %v", err)
		return
	}

	if len(orders) == 0 {
		log.Println("[OrderReconciler] No stale orders found.")
		return
	}

	log.Printf("[OrderReconciler] Found %d stale orders to process.", len(orders))

	// Process orders in batches of ORDER_IN_ONE_CALL
	for i := 0; i < len(orders); i += ORDER_IN_ONE_CALL {
		end := i + ORDER_IN_ONE_CALL
		if end > len(orders) {
			end = len(orders)
		}

		// Get batch of orders
		batch := orders[i:end]
		orderIDs := make([]string, 0, len(batch))
		for _, order := range batch {
			orderIDs = append(orderIDs, order.ID)
		}

		log.Printf("[OrderReconciler] Processing batch of %d orders: %v", len(orderIDs), orderIDs)
		orderReservationCheckingReq := &product_v1.GetOrdersReservationStatusRequest{
			OrderIds: orderIDs,
		}

		go func(req *product_v1.GetOrdersReservationStatusRequest) {
			r.HandleUnreservationOrders(ctx, req)
		}(orderReservationCheckingReq)
	}
}

func (r *OrderReconciler) HandleUnreservationOrders(ctx context.Context, getStatusReq *product_v1.GetOrdersReservationStatusRequest) {
	ordersReservationStatus, err := r.productAdapter.GetOrdersReservationStatus(ctx, getStatusReq)
	if err != nil {
		log.Printf("[OrderReconciler] Error checking order reservation status: %v", err)
		return
	}

	for _, status := range ordersReservationStatus.Orders {
		log.Printf("[OrderReconciler] Order ID: %s, Status: %s", status.OrderId, status.Status)
		// Handle each order status
		if !status.Founded {
			log.Printf("[OrderReconciler] Order ID: %s not found. Updating status to CANCELLED.", status.OrderId)
			// Update the order status to CANCELLED
			_, err := r.orderRepo.UpdateOrderStatus(ctx, status.OrderId, sqlc.OrderStatusCANCELED)
			if err != nil {
				log.Printf("[OrderReconciler] Error updating order status: %v", err)
			}
			continue
		}

		switch status.GetStatus() {
		case product_v1.GetOrderReservationStatusResponse_UNRESERVED.String():
			log.Printf("[OrderReconciler] Order ID: %s is unreserved. No action needed.", status.OrderId)
			// Update the order status to CANCELED
			_, err := r.orderRepo.UpdateOrderStatus(ctx, status.OrderId, sqlc.OrderStatusCANCELED)
			if err != nil {
				log.Printf("[OrderReconciler] Error updating order status to UNRESERVED: %v", err)
			}
		case product_v1.GetOrderReservationStatusResponse_RESERVED.String():
			log.Printf("[OrderReconciler] Order ID: %s is still reserved. Unreserving...", status.OrderId)

			// Prepare unreserve request
			unreserveReq := &product_v1.UnreserveOrdersRequest{
				Orders: []*product_v1.UnreserveOrder{
					{
						OrderId:  status.OrderId,
						ShopId:   status.ShopId,
						Products: []*product_v1.UnreserveProduct{},
					},
				},
			}

			order, err := r.orderRepo.GetOrderByID(ctx, status.OrderId)

			if err != nil {
				log.Printf("[OrderReconciler] Error fetching order details for unreservation: %v", err)
				continue
			}

			// Populate products to unreserve
			for _, item := range order.Items {
				unreserveReq.Orders[0].Products = append(unreserveReq.Orders[0].Products, &product_v1.UnreserveProduct{
					ProductId: item.ProductID,
					Quantity:  int32(item.Quantity),
				})
			}

			// Call product service to unreserve
			resp, err := r.productAdapter.UnreserveOrders(ctx, unreserveReq)
			if err != nil {
				log.Printf("[OrderReconciler] Error calling UnreserveOrders for order %s: %v", status.OrderId, err)
				continue
			}

			// Check result and update order status
			if len(resp.Results) > 0 && resp.Results[0].Success {
				log.Printf("[OrderReconciler] Successfully unreserved order %s", status.OrderId)
				_, err := r.orderRepo.UpdateOrderStatus(ctx, status.OrderId, sqlc.OrderStatusCANCELED)
				if err != nil {
					log.Printf("[OrderReconciler] Error updating order status to CANCELED: %v", err)
				}
			} else {
				log.Printf("[OrderReconciler] Failed to unreserve order %s", status.OrderId)
			}

		default:
			log.Printf("[OrderReconciler] Unknown status for order ID: %s. No action needed.", status.OrderId)
		}
	}
}
