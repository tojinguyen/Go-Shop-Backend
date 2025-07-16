package apperror

// ErrorCode định nghĩa kiểu cho các mã lỗi chuẩn.
type ErrorCode string

// Các hằng số cho mã lỗi. Giúp đảm bảo tính nhất quán và tránh "magic strings".
const (
	// Lỗi chung
	CodeInternal          ErrorCode = "INTERNAL_SERVER_ERROR"
	CodeNotFound          ErrorCode = "NOT_FOUND"
	CodeValidation        ErrorCode = "VALIDATION_ERROR"
	CodeConflict          ErrorCode = "CONFLICT_ERROR"
	CodeUnauthorized      ErrorCode = "UNAUTHORIZED"
	CodeForbidden         ErrorCode = "FORBIDDEN"
	CodeDependencyFailure ErrorCode = "DEPENDENCY_FAILURE" // Lỗi từ service bên ngoài
	CodeNotImplemented    ErrorCode = "NOT_IMPLEMENTED"

	// Lỗi liên quan đến Database
	CodeDatabaseError ErrorCode = "DATABASE_ERROR"

	// Lỗi liên quan đến chuyển đổi dữ liệu
	CodeConversionError ErrorCode = "CONVERSION_ERROR"

	// Lỗi nghiệp vụ cụ thể (có thể mở rộng)
	CodeProductNotFound    ErrorCode = "PRODUCT_NOT_FOUND"
	CodeShopNotFound       ErrorCode = "SHOP_NOT_FOUND"
	CodeUserNotFound       ErrorCode = "USER_NOT_FOUND"
	CodeEmailAlreadyExists ErrorCode = "EMAIL_ALREADY_EXISTS"
	CodeInsufficientStock  ErrorCode = "INSUFFICIENT_STOCK"
	CodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	CodeTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	CodeTokenInvalid       ErrorCode = "TOKEN_INVALID"
)
