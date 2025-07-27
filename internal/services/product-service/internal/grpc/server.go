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
	products, err := s.productRepo.GetByIDs(ctx, req.ProductIds)
	if err != nil {
		log.Printf("Error retrieving products by IDs: %v", err)
		return &product_v1.GetProductsInfoResponse{Valid: false}, err
	}

	if len(products) != len(req.ProductIds) {
		log.Printf("Mismatch count: requested %d, found %d", len(req.ProductIds), len(products))
		return &product_v1.GetProductsInfoResponse{Valid: false}, nil // Some products not found
	}

	var productInfos []*product_v1.ProductInfo
	for _, p := range products {
		productInfos = append(productInfos, &product_v1.ProductInfo{
			Id:       p.ID().String(),
			ShopId:   p.ShopID().String(),
			Price:    int32(p.Price().GetAmount()), // Consider using float or a proper money type
			Currency: p.Price().GetCurrency(),
			Quantity: int32(p.Quantity()),
		})
	}

	return &product_v1.GetProductsInfoResponse{Valid: true, Products: productInfos}, nil
}

func (s *Server) ReserveProducts(ctx context.Context, req *product_v1.ReserveProductsRequest) (*product_v1.ReserveProductsResponse, error) {
	statuses, err := s.productRepo.ReserveStock(ctx, req.GetProducts())
	if err != nil {
		log.Printf("Error during stock reservation for order %s: %v", req.OrderId, err)
		return &product_v1.ReserveProductsResponse{Success: false}, err
	}

	// Check if all reservations were successful
	allSuccess := true
	for _, status := range statuses {
		if !status.Success {
			allSuccess = false
			break
		}
	}

	return &product_v1.ReserveProductsResponse{
		Success:         allSuccess,
		ProductStatuses: statuses,
	}, nil
}

func (s *Server) UnreserveProducts(ctx context.Context, req *product_v1.UnreserveProductsRequest) (*product_v1.UnreserveProductsResponse, error) {
	log.Println("[ProductService] UnreserveProducts called, but not implemented yet.")
	return nil, nil
}
