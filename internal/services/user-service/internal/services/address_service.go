package services

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/converter"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
)

type AddressService struct {
	container *container.ServiceContainer
}

func NewAddressService(container *container.ServiceContainer) *AddressService {
	return &AddressService{
		container: container,
	}
}

func (s *AddressService) CreateAddress(ctx *gin.Context, userID string, req dto.CreateAddressRequest) (*dto.AddressResponse, error) {
	// Tạo parameters cho sqlc
	params := sqlc.CreateAddressParams{
		UserID:    converter.StringToPgUUID(userID),
		IsDefault: converter.BoolToPgBool(req.IsDefault),
		Street:    req.Street,
		Ward:      converter.StringToPgText(req.Ward),
		District:  converter.StringToPgText(req.District),
		City:      converter.StringToPgText(req.City),
		Country:   converter.StringToPgText(&req.Country),
		Lat:       converter.Float64ToPgFloat8(req.Lat),
		Long:      converter.Float64ToPgFloat8(req.Long),
	}

	// Tạo address trong database
	address, err := s.container.GetAddressRepo().CreateAddress(ctx, params)
	if err != nil {
		return nil, err
	}

	// Nếu address này được đặt làm default, cần reset tất cả address khác của user thành false
	if req.IsDefault {
		setDefaultParams := sqlc.SetDefaultAddressParams{
			UserID: converter.StringToPgUUID(userID),
			ID:     converter.StringToPgUUID(address.ID),
		}
		// Update để đảm bảo chỉ address này là default
		_, err = s.container.GetAddressRepo().SetDefaultAddress(ctx, setDefaultParams)
		if err != nil {
			// Log error nhưng không fail request vì address đã được tạo thành công
		}
	}

	// Convert domain model sang DTO response
	response := &dto.AddressResponse{
		ID:        address.ID,
		UserID:    address.UserID,
		IsDefault: address.IsDefault,
		Street:    address.Street,
		Ward:      address.Ward,
		District:  address.District,
		City:      address.City,
		Country:   address.Country,
		Lat:       address.Lat,
		Long:      address.Long,
		DeletedAt: address.DeletedAt,
		CreatedAt: address.CreatedAt,
		UpdatedAt: address.UpdatedAt,
	}

	return response, nil
}

func (s *AddressService) GetAddressByID(ctx *gin.Context, addressID string) (*dto.AddressResponse, error) {
	// Lấy address từ repository
	address, err := s.container.GetAddressRepo().GetAddressByID(ctx, addressID)
	if err != nil {
		log.Printf("Error retrieving address by ID %s: %v", addressID, err)
		return nil, err
	}

	// Chuyển đổi sang DTO response
	response := &dto.AddressResponse{
		ID:        address.ID,
		UserID:    address.UserID,
		IsDefault: address.IsDefault,
		Street:    address.Street,
		Ward:      address.Ward,
		District:  address.District,
		City:      address.City,
		Country:   address.Country,
		Lat:       address.Lat,
		Long:      address.Long,
		DeletedAt: address.DeletedAt,
		CreatedAt: address.CreatedAt,
		UpdatedAt: address.UpdatedAt,
	}

	return response, nil
}

func (s *AddressService) GetAddressesByUserID(ctx *gin.Context, userID string) (*dto.AddressListResponse, error) {
	// Lấy danh sách địa chỉ từ repository
	addresses, err := s.container.GetAddressRepo().GetAddressesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Chuyển đổi sang DTO response
	var addressResponses []dto.AddressResponse
	for _, address := range addresses {
		addressResponses = append(addressResponses, dto.AddressResponse{
			ID:        address.ID,
			UserID:    address.UserID,
			IsDefault: address.IsDefault,
			Street:    address.Street,
			Ward:      address.Ward,
			District:  address.District,
			City:      address.City,
			Country:   address.Country,
			Lat:       address.Lat,
			Long:      address.Long,
			DeletedAt: address.DeletedAt,
			CreatedAt: address.CreatedAt,
			UpdatedAt: address.UpdatedAt,
		})
	}

	return &dto.AddressListResponse{
		Addresses: addressResponses,
		Total:     len(addressResponses),
	}, nil
}

func (s *AddressService) UpdateAddress(ctx *gin.Context, userID string, addressID string, req dto.UpdateAddressRequest) (*dto.AddressResponse, error) {
	// Need to check if address is owned by user
	checkAddress, err := s.container.GetAddressRepo().GetAddressByID(ctx, addressID)
	if err != nil {
		return nil, err
	}
	if checkAddress.UserID != userID {
		return nil, fmt.Errorf("address does not belong to user")
	}

	// Tạo parameters cho sqlc
	params := sqlc.UpdateAddressParams{
		ID:        converter.StringToPgUUID(addressID),
		IsDefault: converter.BoolToPgBool(req.IsDefault),
		Street:    req.Street,
		Ward:      converter.StringToPgText(req.Ward),
		District:  converter.StringToPgText(req.District),
		City:      converter.StringToPgText(req.City),
		Country:   converter.StringToPgText(&req.Country),
		Lat:       converter.Float64ToPgFloat8(req.Lat),
		Long:      converter.Float64ToPgFloat8(req.Long),
	}

	// Cập nhật address trong database
	address, err := s.container.GetAddressRepo().UpdateAddress(ctx, params)
	if err != nil {
		return nil, err
	}

	// Chuyển đổi sang DTO response
	response := &dto.AddressResponse{
		ID:        address.ID,
		UserID:    address.UserID,
		IsDefault: address.IsDefault,
		Street:    address.Street,
		Ward:      address.Ward,
		District:  address.District,
		City:      address.City,
		Country:   address.Country,
		Lat:       address.Lat,
		Long:      address.Long,
		DeletedAt: address.DeletedAt,
		CreatedAt: address.CreatedAt,
		UpdatedAt: address.UpdatedAt,
	}

	return response, nil
}

func (s *AddressService) DeleteAddress(ctx *gin.Context, userID string, addressID string) error {
	// Kiểm tra xem địa chỉ có thuộc về người dùng không
	checkAddress, err := s.container.GetAddressRepo().GetAddressByID(ctx, addressID)
	if err != nil {
		return err
	}
	if checkAddress.UserID != userID {
		return fmt.Errorf("address does not belong to user")
	}

	// Xoá địa chỉ trong database
	err = s.container.GetAddressRepo().DeleteAddress(ctx, addressID)
	if err != nil {
		return err
	}

	return nil
}

func (s *AddressService) SetDefaultAddress(ctx *gin.Context, userID string, addressID string) (*dto.AddressResponse, error) {
	// Kiểm tra xem địa chỉ có thuộc về người dùng không
	checkAddress, err := s.container.GetAddressRepo().GetAddressByID(ctx, addressID)
	if err != nil {
		return nil, err
	}
	if checkAddress.UserID != userID {
		return nil, fmt.Errorf("address does not belong to user")
	}

	// Tạo parameters cho sqlc
	params := sqlc.SetDefaultAddressParams{
		UserID: converter.StringToPgUUID(userID),
		ID:     converter.StringToPgUUID(addressID),
	}

	// Cập nhật địa chỉ thành mặc định
	address, err := s.container.GetAddressRepo().SetDefaultAddress(ctx, params)
	if err != nil {
		return nil, err
	}

	// Chuyển đổi sang DTO response
	response := &dto.AddressResponse{
		ID:        address.ID,
		UserID:    address.UserID,
		IsDefault: address.IsDefault,
		Street:    address.Street,
		Ward:      address.Ward,
		District:  address.District,
		City:      address.City,
		Country:   address.Country,
		Lat:       address.Lat,
		Long:      address.Long,
		DeletedAt: address.DeletedAt,
		CreatedAt: address.CreatedAt,
		UpdatedAt: address.UpdatedAt,
	}

	return response, nil
}
