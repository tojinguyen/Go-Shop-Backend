package services_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/toji-dev/go-shop/internal/pkg/constant"
	user_constant "github.com/toji-dev/go-shop/internal/services/user-service/internal/constant"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	repo_mocks "github.com/toji-dev/go-shop/internal/services/user-service/internal/repository/mocks"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/services"
)

func TestUserService_CreateProfile(t *testing.T) {
	testCases := []struct {
		name          string
		request       dto.CreateUserProfileRequest
		setupMocks    func(mockRepo *repo_mocks.UserProfileRepository)
		expectedID    string
		expectError   bool
		expectedError string
	}{
		{
			name:    "Success - User profile created successfully",
			request: dto.CreateUserProfileRequest{FullName: "Toai Nguyen"},
			setupMocks: func(mockRepo *repo_mocks.UserProfileRepository) {
				createdUserProfile := &domain.UserProfile{
					UserID:           "550e8400-e29b-41d4-a716-446655440000",
					Email:            "H0G0Y@example.com",
					FullName:         "Toai Nguyen",
					Birthday:         "1990-01-01",
					Phone:            "1234567890",
					Role:             string(user_constant.UserRoleCustomer),
					BannedAt:         "",
					AvatarURL:        "https://example.com/avatar.jpg",
					DefaultAddressID: "",
					Gender:           string(user_constant.UserGenderMale),
					CreatedAt:        "2023-10-01T00:00:00Z",
					UpdatedAt:        "2023-10-01T00:00:00Z",
				}

				mockRepo.EXPECT().CreateUserProfile(mock.Anything, mock.AnythingOfType("sqlc.CreateUserProfileParams")).Return(createdUserProfile, nil).Once()
			},
			expectedID:  "550e8400-e29b-41d4-a716-446655440000",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(repo_mocks.UserProfileRepository)
			tc.setupMocks(mockRepo)

			userService := services.NewUserService(mockRepo)

			ctx := createTestGinContext()
			ctx.Set(constant.ContextKeyUserID, "550e8400-e29b-41d4-a716-446655440000")
			ctx.Set(constant.ContextKeyUserRole, string(user_constant.UserRoleCustomer))
			ctx.Set(constant.ContextKeyUserEmail, "H0G0Y@example.com")

			result, err := userService.CreateProfile(ctx, tc.request)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if err.Error() != tc.expectedError {
					t.Errorf("expected error %s but got %s", tc.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				} else if result.UserID != tc.expectedID {
					t.Errorf("expected ID %s but got %s", tc.expectedID, result.UserID)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
