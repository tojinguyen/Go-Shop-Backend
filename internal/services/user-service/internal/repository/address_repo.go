package repository

import (
	"context"
	"log"

	"github.com/toji-dev/go-shop/internal/pkg/converter"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/domain"
)

type AddressRepository interface {
	CreateAddress(ctx context.Context, params sqlc.CreateAddressParams) (*domain.Address, error)
	GetAddressByID(ctx context.Context, addressID string) (*domain.Address, error)
	UpdateAddress(ctx context.Context, params sqlc.UpdateAddressParams) (*domain.Address, error)
	DeleteAddress(ctx context.Context, addressID string) error
	GetAddressesByUserID(ctx context.Context, userID string) ([]domain.Address, error)
	GetDefaultAddressByUserID(ctx context.Context, userID string) (*domain.Address, error)
	SetDefaultAddress(ctx context.Context, params sqlc.SetDefaultAddressParams) (*domain.Address, error)
}

type addressRepository struct {
	db      *postgresql_infra.PostgreSQLService
	queries *sqlc.Queries
}

func NewAddressRepository(db *postgresql_infra.PostgreSQLService) AddressRepository {
	return &addressRepository{
		db:      db,
		queries: sqlc.New(db.GetPool()),
	}
}

func (r *addressRepository) CreateAddress(ctx context.Context, params sqlc.CreateAddressParams) (*domain.Address, error) {
	address, err := r.queries.CreateAddress(ctx, params)
	if err != nil {
		return nil, err
	}

	return &domain.Address{
		ID:        address.ID.String(),
		UserID:    address.UserID.String(),
		IsDefault: address.IsDefault.Bool,
		Street:    address.Street,
		Ward:      *converter.PgTextToStringPtr(address.Ward),
		District:  *converter.PgTextToStringPtr(address.District),
		City:      *converter.PgTextToStringPtr(address.City),
		Country:   address.Country.String,
		Lat:       *converter.PgFloat8ToFloat64Ptr(address.Lat),
		Long:      *converter.PgFloat8ToFloat64Ptr(address.Long),
		CreatedAt: address.CreatedAt.Time,
		UpdatedAt: address.UpdatedAt.Time,
	}, nil
}

func (r *addressRepository) GetAddressByID(ctx context.Context, addressID string) (*domain.Address, error) {
	log.Printf("Fetching address with ID: %s", addressID)
	address, err := r.queries.GetAddressById(ctx, converter.StringToPgUUID(addressID))
	if err != nil {
		log.Printf("Error fetching address: %v", err)
		return nil, err
	}

	addressDomain := &domain.Address{
		ID:        address.ID.String(),
		UserID:    address.UserID.String(),
		IsDefault: address.IsDefault.Bool,
		Street:    address.Street,
		Ward:      *converter.PgTextToStringPtr(address.Ward),
		District:  *converter.PgTextToStringPtr(address.District),
		City:      *converter.PgTextToStringPtr(address.City),
		Country:   address.Country.String,
		Lat:       *converter.PgFloat8ToFloat64Ptr(address.Lat),
		Long:      *converter.PgFloat8ToFloat64Ptr(address.Long),
		CreatedAt: address.CreatedAt.Time,
		UpdatedAt: address.UpdatedAt.Time,
	}

	if address.DeletedAt.Valid {
		addressDomain.DeletedAt = address.DeletedAt.Time
	}

	return addressDomain, nil
}

func (r *addressRepository) UpdateAddress(ctx context.Context, params sqlc.UpdateAddressParams) (*domain.Address, error) {
	address, err := r.queries.UpdateAddress(ctx, params)
	if err != nil {
		return nil, err
	}

	addressRes := &domain.Address{
		ID:        address.ID.String(),
		UserID:    address.UserID.String(),
		IsDefault: address.IsDefault.Bool,
		Street:    address.Street,
		Ward:      *converter.PgTextToStringPtr(address.Ward),
		District:  *converter.PgTextToStringPtr(address.District),
		City:      *converter.PgTextToStringPtr(address.City),
		Country:   address.Country.String,
		Lat:       *converter.PgFloat8ToFloat64Ptr(address.Lat),
		Long:      *converter.PgFloat8ToFloat64Ptr(address.Long),
		CreatedAt: address.CreatedAt.Time,
		UpdatedAt: address.UpdatedAt.Time,
	}

	if address.DeletedAt.Valid {
		addressRes.DeletedAt = address.DeletedAt.Time
	}

	return addressRes, nil
}

func (r *addressRepository) DeleteAddress(ctx context.Context, addressID string) error {
	err := r.queries.DeleteAddress(ctx, converter.StringToPgUUID(addressID))
	if err != nil {
		return err
	}
	return nil
}

func (r *addressRepository) GetAddressesByUserID(ctx context.Context, userID string) ([]domain.Address, error) {
	addresses, err := r.queries.GetAddressesByUserId(ctx, converter.StringToPgUUID(userID))
	if err != nil {
		return nil, err
	}

	var result []domain.Address
	for _, address := range addresses {
		result = append(result, domain.Address{
			ID:        address.ID.String(),
			UserID:    address.UserID.String(),
			IsDefault: address.IsDefault.Bool,
			Street:    address.Street,
			Ward:      *converter.PgTextToStringPtr(address.Ward),
			District:  *converter.PgTextToStringPtr(address.District),
			City:      *converter.PgTextToStringPtr(address.City),
			Country:   address.Country.String,
			Lat:       *converter.PgFloat8ToFloat64Ptr(address.Lat),
			Long:      *converter.PgFloat8ToFloat64Ptr(address.Long),
			CreatedAt: address.CreatedAt.Time,
			UpdatedAt: address.UpdatedAt.Time,
		})

		if address.DeletedAt.Valid {
			result[len(result)-1].DeletedAt = address.DeletedAt.Time
		}
	}

	return result, nil
}

func (r *addressRepository) GetDefaultAddressByUserID(ctx context.Context, userID string) (*domain.Address, error) {
	address, err := r.queries.GetDefaultAddressByUserId(ctx, converter.StringToPgUUID(userID))
	if err != nil {
		return nil, err
	}

	addressRes := &domain.Address{
		ID:        address.ID.String(),
		UserID:    address.UserID.String(),
		IsDefault: address.IsDefault.Bool,
		Street:    address.Street,
		Ward:      *converter.PgTextToStringPtr(address.Ward),
		District:  *converter.PgTextToStringPtr(address.District),
		City:      *converter.PgTextToStringPtr(address.City),
		Country:   address.Country.String,
		Lat:       *converter.PgFloat8ToFloat64Ptr(address.Lat),
		Long:      *converter.PgFloat8ToFloat64Ptr(address.Long),
		CreatedAt: address.CreatedAt.Time,
		UpdatedAt: address.UpdatedAt.Time,
	}

	if address.DeletedAt.Valid {
		addressRes.DeletedAt = address.DeletedAt.Time
	}

	return addressRes, nil
}

func (r *addressRepository) SetDefaultAddress(ctx context.Context, params sqlc.SetDefaultAddressParams) (*domain.Address, error) {
	err := r.queries.SetDefaultAddress(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.GetAddressByID(ctx, params.ID.String())
}
