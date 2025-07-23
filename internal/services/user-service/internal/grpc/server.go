package grpc

import (
	"context"

	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/repository"
	user_v1 "github.com/toji-dev/go-shop/proto/gen/go/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	user_v1.UnimplementedUserServiceServer
	addressRepo repository.AddressRepository
}

func NewUserGRPCServer(addressRepo repository.AddressRepository) *Server {
	return &Server{
		addressRepo: addressRepo,
	}
}

func (s *Server) GetAddressById(ctx context.Context, req *user_v1.GetAddressRequest) (*user_v1.GetAddressResponse, error) {
	address, err := s.addressRepo.GetAddressByID(ctx, req.AddressId)
	if err != nil {
		return nil, err
	}

	if address == nil {
		return nil, apperror.NewNotFound("Address", req.AddressId)
	}

	addressProto := &user_v1.Address{
		Id:        address.ID,
		UserId:    address.UserID,
		IsDefault: address.IsDefault,
		Street:    address.Street,
		Ward:      address.Ward,
		District:  address.District,
		City:      address.City,
		Country:   address.Country,
		Lat:       address.Lat,
		Long:      address.Long,
		DeletedAt: timestamppb.New(address.DeletedAt),
		CreatedAt: timestamppb.New(address.CreatedAt),
		UpdatedAt: timestamppb.New(address.UpdatedAt),
	}

	return &user_v1.GetAddressResponse{
		Address: addressProto,
	}, nil
}
