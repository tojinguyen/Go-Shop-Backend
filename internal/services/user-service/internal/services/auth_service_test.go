package services_test

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	email_mocks "github.com/toji-dev/go-shop/internal/pkg/email/mocks"
	redis_mocks "github.com/toji-dev/go-shop/internal/pkg/infra/redis-infra/mocks"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"
	jwt_mocks "github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/jwt/mocks"
	repo_mocks "github.com/toji-dev/go-shop/internal/services/user-service/internal/repository/mocks"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/services"
)

func createTestGinContext() *gin.Context {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)
	return ctx
}

// =================================================================
// Test cho hàm Register
// =================================================================
func TestAuthService_Register(t *testing.T) {
	// Định nghĩa các kịch bản test trong một slice (table-driven test)
	testCases := []struct {
		name          string
		request       dto.RegisterRequest
		setupMocks    func(mockRepo *repo_mocks.UserAccountRepository)
		expectedID    string
		expectError   bool
		expectedError string
	}{
		{
			name: "Success - User registered successfully",
			request: dto.RegisterRequest{
				Email:    "newuser@example.com",
				Password: "Password123!",
			},
			setupMocks: func(mockRepo *repo_mocks.UserAccountRepository) {
				// 1. Giả lập việc kiểm tra user chưa tồn tại
				mockRepo.On("CheckUserExistsByEmail", mock.Anything, "newuser@example.com").Return(false, nil).Once()

				// 2. Giả lập việc tạo user thành công và trả về user mới
				createdUser := &domain.UserAccount{Id: "new-user-id", Email: "newuser@example.com"}
				// mock.AnythingOfType dùng để khớp với bất kỳ giá trị nào của kiểu dữ liệu đó
				mockRepo.On("CreateUserAccount", mock.Anything, mock.AnythingOfType("sqlc.CreateUserAccountParams")).Return(createdUser, nil).Once()
			},
			expectedID:  "new-user-id",
			expectError: false,
		},
		{
			name: "Error - User already exists",
			request: dto.RegisterRequest{
				Email:    "existing@example.com",
				Password: "Password123!",
			},
			setupMocks: func(mockRepo *repo_mocks.UserAccountRepository) {
				// 1. Giả lập việc kiểm tra user đã tồn tại
				mockRepo.On("CheckUserExistsByEmail", mock.Anything, "existing@example.com").Return(true, nil).Once()
				// Hàm CreateUserAccount sẽ không được gọi, nên không cần mock
			},
			expectError:   true,
			expectedError: "user with email existing@example.com already exists",
		},
		{
			name: "Error - Database error on checking existence",
			request: dto.RegisterRequest{
				Email:    "dberror@example.com",
				Password: "Password123!",
			},
			setupMocks: func(mockRepo *repo_mocks.UserAccountRepository) {
				// 1. Giả lập DB trả về lỗi khi kiểm tra user
				mockRepo.On("CheckUserExistsByEmail", mock.Anything, "dberror@example.com").Return(false, errors.New("db connection error")).Once()
			},
			expectError:   true,
			expectedError: "failed to check user existence",
		},
	}

	// Chạy vòng lặp qua tất cả các test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// --- ARRANGE ---
			// Khởi tạo các mock objects cần thiết cho AuthService
			mockUserRepo := new(repo_mocks.UserAccountRepository)
			mockJwtService := new(jwt_mocks.JwtService)
			mockRedisService := new(redis_mocks.RedisServiceInterface)
			mockEmailService := new(email_mocks.EmailService)

			// Áp dụng thiết lập mock cho kịch bản hiện tại
			tc.setupMocks(mockUserRepo)

			// Tạo instance của AuthService với các dependency đã được mock
			authService := services.NewAuthService(
				mockUserRepo,
				mockJwtService,
				mockRedisService,
				mockEmailService,
				&config.Config{}, // Config rỗng vì hàm Register không dùng đến
			)

			// --- ACT ---
			response, err := authService.Register(createTestGinContext(), tc.request)

			// --- ASSERT ---
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedID, response.UserID)
			}

			// Xác minh rằng các hàm mock đã được gọi đúng như kỳ vọng
			mockUserRepo.AssertExpectations(t)
		})
	}
}
