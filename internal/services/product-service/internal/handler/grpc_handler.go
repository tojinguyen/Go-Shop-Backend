package handler

import (
	"context"
	"log"

	"github.com/toji-dev/go-shop/internal/services/product-service/internal/service"
	pb "github.com/toji-dev/go-shop/proto/product"
)

type GRPCHandler struct {
	pb.UnimplementedProductServiceServer
	productService *service.ProductService
}

func NewGRPCHandler(productService *service.ProductService) *GRPCHandler {
	return &GRPCHandler{
		productService: productService,
	}
}

func (h *GRPCHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	log.Printf("Received gRPC CreateProduct request for shop %s", req.ShopId)

	_, err := h.productService.CreateProduct(ctx, req)

	if err != nil {
		log.Printf("Error creating product: %v", err)
		return nil, err
	}

	return &pb.CreateProductResponse{}, nil
}
