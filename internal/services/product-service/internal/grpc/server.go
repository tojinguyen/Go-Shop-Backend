package grpc

import (
	"context"
	"log"

	"github.com/toji-dev/go-shop/internal/services/product-service/internal/repository"
	product_v1 "github.com/toji-dev/go-shop/proto/gen/go/product/v1"
)

type Server struct {
	product_v1.UnimplementedProductServiceServer
	productRepo repository.ProductRepository
}

func NewProductGRPCServer(productRepo repository.ProductRepository) *Server {
	return &Server{
		productRepo: productRepo,
	}
}

func (s *Server) GetProductInfo(ctx context.Context, req *product_v1.GetProductInfoRequest) (*product_v1.GetProductInfoResponse, error) {
	product, err := s.productRepo.GetByID(ctx, req.ProductId)
	if err != nil {
		log.Printf("Error retrieving product with ID %s: %v", req.ProductId, err)
		return &product_v1.GetProductInfoResponse{
			Exists:  false,
			Product: nil,
		}, err
	}

	if product == nil || product.DeletedAt() != nil {
		log.Printf("Product with ID %s not found", req.ProductId)
		return &product_v1.GetProductInfoResponse{
			Exists:  false,
			Product: nil,
		}, nil
	}

	productInfo := &product_v1.ProductInfo{
		Id:       product.ID().String(),
		ShopId:   product.ShopID().String(),
		Price:    int32(product.Price().GetAmount()),
		Currency: product.Price().GetCurrency(),
		Quantity: int32(product.Quantity()),
	}

	return &product_v1.GetProductInfoResponse{
		Exists:  product != nil,
		Product: productInfo,
	}, nil
}

func (s *Server) GetProductsInfo(ctx context.Context, req *product_v1.GetProductsInfoRequest) (*product_v1.GetProductsInfoResponse, error) {
	return nil, nil
}
