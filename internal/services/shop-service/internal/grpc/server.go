package grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	promotionRepo "github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository/promotion"
	shopRepo "github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository/shop"
	shop_v1 "github.com/toji-dev/go-shop/proto/gen/go/shop/v1"
)

type Server struct {
	shop_v1.UnimplementedShopServiceServer
	shopRepo      shopRepo.ShopRepository
	promotionRepo promotionRepo.PromotionRepository
}

func NewShopGRPCServer(repo shopRepo.ShopRepository, promotionRepo promotionRepo.PromotionRepository) *Server {
	return &Server{
		shopRepo:      repo,
		promotionRepo: promotionRepo,
	}
}

func (s *Server) CheckShopOwnership(ctx context.Context, req *shop_v1.CheckShopOwnershipRequest) (*shop_v1.CheckShopOwnershipResponse, error) {
	log.Printf("Received CheckShopOwnership request for ShopID: %s, UserID: %s", req.GetShopId(), req.GetUserId())

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Printf("Invalid UserID format: %s", req.UserId)
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	shopID, err := uuid.Parse(req.ShopId)
	if err != nil {
		log.Printf("Invalid ShopID format: %s", req.ShopId)
		return nil, fmt.Errorf("invalid shop ID format: %w", err)
	}

	shopInfo, err := s.shopRepo.GetShopByID(ctx, shopID.String())
	if err != nil {
		log.Printf("Error retrieving shop with ID %s: %v", shopID, err)
		return nil, fmt.Errorf("error retrieving shop: %w", err)
	}

	if shopInfo == nil {
		log.Printf("Shop with ID %s not found", shopID)
		return nil, fmt.Errorf("shop not found")
	}

	isOwner := shopInfo.OwnerID == userID
	log.Printf("User %s ownership status for Shop %s: %t", userID, shopID, isOwner)

	return &shop_v1.CheckShopOwnershipResponse{
		IsOwner: isOwner,
	}, nil
}

func (s *Server) CheckShopExists(ctx context.Context, req *shop_v1.CheckShopExistsRequest) (*shop_v1.CheckShopExistsResponse, error) {
	log.Printf("Received CheckShopExists request for ShopID: %s", req.GetShopId())

	shopID, err := uuid.Parse(req.ShopId)
	if err != nil {
		log.Printf("Invalid ShopID format: %s", req.ShopId)
		return nil, fmt.Errorf("invalid shop ID format: %w", err)
	}

	shopInfo, err := s.shopRepo.GetShopByID(ctx, shopID.String())
	if err != nil {
		log.Printf("Error checking existence of shop with ID %s: %v", shopID, err)
		return nil, fmt.Errorf("error checking shop existence: %w", err)
	}

	exists := shopInfo != nil

	return &shop_v1.CheckShopExistsResponse{
		Exists: exists,
	}, nil
}

func (s *Server) CalculatePromotion(ctx context.Context, req *shop_v1.CalculatePromotionRequest) (*shop_v1.CalculatePromotionResponse, error) {
	log.Printf("Received CalculatePromotion request for ShopID: %s, UserID: %s, PromotionCode: %s, TotalAmount: %d",
		req.GetShopId(), req.GetUserId(), req.GetPromotionCode(), req.GetTotalAmount())

	promotion, err := s.promotionRepo.GetByID(ctx, req.GetPromotionCode())
	if err != nil {
		log.Printf("Error retrieving promotions for ShopID %s: %v", req.GetShopId(), err)
		return nil, fmt.Errorf("error retrieving promotions: %s", err)
	}

	discount, err := promotion.CalculateDiscount(float64(req.GetTotalAmount()))
	if err != nil {
		log.Printf("Error calculating discount for promotion %s: %v", req.GetPromotionCode(), err)
		return nil, fmt.Errorf("error calculating discount: %s", err)
	}

	return &shop_v1.CalculatePromotionResponse{
		Eligible: true,
		Discount: float32(discount),
	}, nil
}
