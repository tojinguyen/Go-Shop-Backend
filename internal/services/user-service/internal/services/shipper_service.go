package services

import (
	"context"

	"github.com/toji-dev/go-shop/internal/services/user-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/repository"
)

type ShipperService struct {
	shipperRepo repository.ShipperRepository
}

func NewShipperService(shipperRepo repository.ShipperRepository) *ShipperService {
	return &ShipperService{
		shipperRepo: shipperRepo,
	}
}

func (s *ShipperService) RegisterShipper(ctx context.Context, userID string, request *dto.ShipperRegisterRequest) (*dto.ShipperResponse, error) {
	shipper, err := s.shipperRepo.CreateShipper(ctx, &domain.Shipper{
		UserID:          userID,
		VehicleType:     request.VehicleType,
		VehicleImageURL: request.VehicleImageURL,
		IdentifyCardURL: request.IdentifyCardURL,
		LicensePlate:    request.LicensePlate,
	})
	if err != nil {
		return nil, err
	}

	return &dto.ShipperResponse{
		UserID:          shipper.UserID,
		VehicleType:     shipper.VehicleType,
		VehicleImageURL: shipper.VehicleImageURL,
		IdentifyCardURL: shipper.IdentifyCardURL,
		LicensePlate:    shipper.LicensePlate,
		CreatedAt:       shipper.CreatedAt,
		UpdatedAt:       shipper.UpdatedAt,
	}, nil
}

func (s *ShipperService) GetShipperProfile(ctx context.Context, userID string) (*dto.ShipperResponse, error) {
	shipper, err := s.shipperRepo.GetShipperByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if shipper == nil {
		return nil, nil
	}

	return &dto.ShipperResponse{
		UserID:          shipper.UserID,
		VehicleType:     shipper.VehicleType,
		VehicleImageURL: shipper.VehicleImageURL,
		IdentifyCardURL: shipper.IdentifyCardURL,
		LicensePlate:    shipper.LicensePlate,
		CreatedAt:       shipper.CreatedAt,
		UpdatedAt:       shipper.UpdatedAt,
	}, nil
}

func (s *ShipperService) UpdateShipperProfile(ctx context.Context, userID string, request *dto.ShipperUpdateRequest) (*dto.ShipperResponse, error) {
	shipper, err := s.shipperRepo.UpdateShipper(ctx, userID, request)
	if err != nil {
		return nil, err
	}

	return &dto.ShipperResponse{
		UserID:          shipper.UserID,
		VehicleType:     shipper.VehicleType,
		VehicleImageURL: shipper.VehicleImageURL,
		IdentifyCardURL: shipper.IdentifyCardURL,
		LicensePlate:    shipper.LicensePlate,
	}, nil
}

func (s *ShipperService) DeleteShipperProfile(ctx context.Context, userID string) error {
	err := s.shipperRepo.DeleteShipper(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}
