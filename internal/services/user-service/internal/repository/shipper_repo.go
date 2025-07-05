package repository

import (
	"context"
	"log"

	"github.com/toji-dev/go-shop/internal/pkg/converter"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
)

type ShipperRepository interface {
	CreateShipper(ctx context.Context, shipper *domain.Shipper) (*domain.Shipper, error)
	GetShipperByUserID(ctx context.Context, userID string) (*domain.Shipper, error)
	UpdateShipper(ctx context.Context, userID string, updateRequest *dto.ShipperUpdateRequest) (*domain.Shipper, error)
	DeleteShipper(ctx context.Context, userID string) error
}

type shipperRepository struct {
	db      *postgresql_infra.PostgreSQLService
	queries *sqlc.Queries
}

func NewShipperRepository(db *postgresql_infra.PostgreSQLService) ShipperRepository {
	return &shipperRepository{
		db:      db,
		queries: sqlc.New(db.GetPool()),
	}
}

// CreateShipper creates a new shipper profile
func (r *shipperRepository) CreateShipper(ctx context.Context, shipper *domain.Shipper) (*domain.Shipper, error) {
	params := sqlc.CreateShipperParams{
		UserID:          converter.StringToPgUUID(shipper.UserID),
		VehicleType:     converter.StringToPgText(shipper.VehicleType),
		VehicleImageUrl: converter.StringToPgText(shipper.VehicleImageURL),
		IdentifyCardUrl: converter.StringToPgText(shipper.IdentifyCardURL),
		LicensePlate:    converter.StringToPgText(shipper.LicensePlate),
	}

	result, err := r.queries.CreateShipper(ctx, params)
	if err != nil {
		return nil, err
	}

	return &domain.Shipper{
		UserID:          converter.PgUUIDToString(result.UserID),
		VehicleType:     converter.PgTextToStringPtr(result.VehicleType),
		VehicleImageURL: converter.PgTextToStringPtr(result.VehicleImageUrl),
		IdentifyCardURL: converter.PgTextToStringPtr(result.IdentifyCardUrl),
		LicensePlate:    converter.PgTextToStringPtr(result.LicensePlate),
	}, nil
}

// GetShipperByUserID retrieves a shipper profile by user ID
func (r *shipperRepository) GetShipperByUserID(ctx context.Context, userID string) (*domain.Shipper, error) {
	result, err := r.queries.GetShipperByUserID(ctx, converter.StringToPgUUID(userID))
	if err != nil {
		return nil, err
	}

	return &domain.Shipper{
		UserID:          converter.PgUUIDToString(result.UserID),
		VehicleType:     converter.PgTextToStringPtr(result.VehicleType),
		VehicleImageURL: converter.PgTextToStringPtr(result.VehicleImageUrl),
		IdentifyCardURL: converter.PgTextToStringPtr(result.IdentifyCardUrl),
		LicensePlate:    converter.PgTextToStringPtr(result.LicensePlate),
	}, nil
}

// UpdateShipper updates an existing shipper profile
func (r *shipperRepository) UpdateShipper(ctx context.Context, userID string, updateRequest *dto.ShipperUpdateRequest) (*domain.Shipper, error) {
	// First get the current shipper to preserve existing values for fields not being updated
	currentShipper, err := r.GetShipperByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Prepare update parameters, using current values if new ones are not provided
	params := sqlc.UpdateShipperByUserIDParams{
		UserID:          converter.StringToPgUUID(userID),
		VehicleType:     converter.StringToPgText(updateRequest.VehicleType),
		VehicleImageUrl: converter.StringToPgText(updateRequest.VehicleImageURL),
		IdentifyCardUrl: converter.StringToPgText(updateRequest.IdentifyCardURL),
		LicensePlate:    converter.StringToPgText(updateRequest.LicensePlate),
	}

	// If fields are nil in update request, use current values
	if updateRequest.VehicleType == nil {
		params.VehicleType = converter.StringToPgText(currentShipper.VehicleType)
	}
	if updateRequest.VehicleImageURL == nil {
		params.VehicleImageUrl = converter.StringToPgText(currentShipper.VehicleImageURL)
	}
	if updateRequest.IdentifyCardURL == nil {
		params.IdentifyCardUrl = converter.StringToPgText(currentShipper.IdentifyCardURL)
	}
	if updateRequest.LicensePlate == nil {
		params.LicensePlate = converter.StringToPgText(currentShipper.LicensePlate)
	}

	result, err := r.queries.UpdateShipperByUserID(ctx, params)
	if err != nil {
		log.Printf("Error updating shipper: %v", err)
		return nil, err
	}

	return &domain.Shipper{
		UserID:          converter.PgUUIDToString(result.UserID),
		VehicleType:     converter.PgTextToStringPtr(result.VehicleType),
		VehicleImageURL: converter.PgTextToStringPtr(result.VehicleImageUrl),
		IdentifyCardURL: converter.PgTextToStringPtr(result.IdentifyCardUrl),
		LicensePlate:    converter.PgTextToStringPtr(result.LicensePlate),
	}, nil
}

// DeleteShipper deletes a shipper profile
func (r *shipperRepository) DeleteShipper(ctx context.Context, userID string) error {
	// First check if shipper exists
	_, err := r.queries.GetShipperByUserID(ctx, converter.StringToPgUUID(userID))
	if err != nil {
		return err
	}

	// Delete the shipper profile
	return r.queries.DeleteShipperByUserID(ctx, converter.StringToPgUUID(userID))
}
