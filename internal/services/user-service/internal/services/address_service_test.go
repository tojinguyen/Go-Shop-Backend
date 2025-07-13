package services_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	time_utils "github.com/toji-dev/go-shop/internal/pkg/time"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	repo_mocks "github.com/toji-dev/go-shop/internal/services/user-service/internal/repository/mocks"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/services"
)

func TestAddressService_CreateAddress(t *testing.T) {
	ctx := createTestGinContext()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	addressID := "a2b7a59a-2d31-4b6a-8c7b-9e8a7b6d5c4d"

	testCases := []struct {
		name          string
		setupMocks    func(mockRepo *repo_mocks.AddressRepository)
		req           dto.CreateAddressRequest
		expectedError string
	}{
		{
			name: "Success - Create default address",
			setupMocks: func(mockRepo *repo_mocks.AddressRepository) {
				mockRepo.EXPECT().CreateAddress(mock.Anything, mock.AnythingOfType("sqlc.CreateAddressParams")).Return(&domain.Address{ID: addressID}, nil).Once()
				mockRepo.EXPECT().SetDefaultAddress(mock.Anything, mock.AnythingOfType("sqlc.SetDefaultAddressParams")).Return(&domain.Address{}, nil).Once()
			},
			req: dto.CreateAddressRequest{
				Street:    "123 Main St",
				City:      "Anytown",
				Country:   "USA",
				IsDefault: true,
			},
			expectedError: "",
		},
		{
			name: "Success - Create non-default address",
			setupMocks: func(mockRepo *repo_mocks.AddressRepository) {
				mockRepo.EXPECT().CreateAddress(mock.Anything, mock.AnythingOfType("sqlc.CreateAddressParams")).Return(&domain.Address{ID: addressID}, nil).Once()
			},
			req: dto.CreateAddressRequest{
				Street:    "456 Oak Ave",
				City:      "Otherville",
				Country:   "USA",
				IsDefault: false,
			},
			expectedError: "",
		},
		{
			name: "Error - Repository fails to create address",
			setupMocks: func(mockRepo *repo_mocks.AddressRepository) {
				mockRepo.EXPECT().CreateAddress(mock.Anything, mock.AnythingOfType("sqlc.CreateAddressParams")).Return(nil, fmt.Errorf("database error")).Once()
			},
			req:           dto.CreateAddressRequest{},
			expectedError: "database error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := repo_mocks.NewAddressRepository(t)
			addressService := services.NewAddressService(mockRepo)
			tc.setupMocks(mockRepo)

			_, err := addressService.CreateAddress(ctx, userID, tc.req)

			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddressService_GetAddressByID(t *testing.T) {
	ctx := createTestGinContext()
	addressID := "a2b7a59a-2d31-4b6a-8c7b-9e8a7b6d5c4d"

	testCases := []struct {
		name          string
		setupMocks    func(mockRepo *repo_mocks.AddressRepository)
		addressID     string
		expectedError string
	}{
		{
			name: "Success - Get address by ID",
			setupMocks: func(mockRepo *repo_mocks.AddressRepository) {
				mockRepo.EXPECT().GetAddressByID(mock.Anything, addressID).Return(&domain.Address{
					ID:        addressID,
					UserID:    "550e8400-e29b-41d4-a716-446655440000",
					Street:    "123 Main St",
					City:      "Anytown",
					Country:   "USA",
					CreatedAt: time_utils.GetUtcTime(),
					UpdatedAt: time_utils.GetUtcTime(),
				}, nil).Once()
			},
			addressID:     addressID,
			expectedError: "",
		},
		{
			name: "Error - Address not found",
			setupMocks: func(mockRepo *repo_mocks.AddressRepository) {
				mockRepo.EXPECT().GetAddressByID(mock.Anything, "non-existent-id").Return(nil, fmt.Errorf("not found")).Once()
			},
			addressID:     "non-existent-id",
			expectedError: "not found",
		},
		{
			name: "Error - Repository returns an error",
			setupMocks: func(mockRepo *repo_mocks.AddressRepository) {
				mockRepo.EXPECT().GetAddressByID(mock.Anything, addressID).Return(nil, fmt.Errorf("database error")).Once()
			},
			addressID:     addressID,
			expectedError: "database error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := repo_mocks.NewAddressRepository(t)
			addressService := services.NewAddressService(mockRepo)
			tc.setupMocks(mockRepo)

			address, err := addressService.GetAddressByID(ctx, tc.addressID)

			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
				assert.Nil(t, address)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, address)
				assert.Equal(t, tc.addressID, address.ID)
			}
		})
	}
}
