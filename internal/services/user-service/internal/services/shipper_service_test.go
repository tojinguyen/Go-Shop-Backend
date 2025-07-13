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

func TestShipperService_RegisterShipper(t *testing.T) {
	ctx := createTestGinContext()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	shipper := &domain.Shipper{
		UserID:          userID,
		VehicleType:     "Yamaha",
		VehicleImageURL: "http://example.com/vehicle.jpg",
		IdentifyCardURL: "http://example.com/identify_card.jpg",
		LicensePlate:    "ABC-1234",
		CreatedAt:       time_utils.GetUtcTime(),
		UpdatedAt:       time_utils.GetUtcTime(),
	}

	testCases := []struct {
		name          string
		setupMocks    func(mockRepo *repo_mocks.ShipperRepository)
		req           *dto.ShipperRegisterRequest
		expectedError string
	}{
		{
			name: "Success - Register shipper",
			setupMocks: func(mockRepo *repo_mocks.ShipperRepository) {
				mockRepo.EXPECT().CreateShipper(mock.Anything, mock.AnythingOfType("*domain.Shipper")).Return(shipper, nil).Once()
			},

			req: &dto.ShipperRegisterRequest{
				VehicleType:     shipper.VehicleType,
				VehicleImageURL: shipper.VehicleImageURL,
				IdentifyCardURL: shipper.IdentifyCardURL,
				LicensePlate:    shipper.LicensePlate,
			},
			expectedError: "",
		},
		{
			name: "Error - Repository fails to create shipper",
			setupMocks: func(mockRepo *repo_mocks.ShipperRepository) {
				mockRepo.EXPECT().CreateShipper(mock.Anything, mock.AnythingOfType("*domain.Shipper")).Return(nil, fmt.Errorf("database error")).Once()
			},
			req: &dto.ShipperRegisterRequest{
				VehicleType:     shipper.VehicleType,
				VehicleImageURL: shipper.VehicleImageURL,
				IdentifyCardURL: shipper.IdentifyCardURL,
				LicensePlate:    shipper.LicensePlate,
			},
			expectedError: "database error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := repo_mocks.NewShipperRepository(t)
			tc.setupMocks(mockRepo)

			service := services.NewShipperService(mockRepo)
			resp, err := service.RegisterShipper(ctx, userID, tc.req)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, shipper.UserID, resp.UserID)
			}
		})
	}
}
