package apperror

import (
	"errors"
	"fmt"
)

// ErrorType định danh loại lỗi để handler có thể map sang HTTP status code.
type ErrorType int

const (
	// Các loại lỗi nghiệp vụ và hệ thống
	TypeNotFound          ErrorType = iota // 404
	TypeValidation                         // 400
	TypeConflict                           // 409 (e.g., duplicate entry)
	TypeUnauthorized                       // 401
	TypeForbidden                          // 403
	TypeInternal                           // 500
	TypeDependencyFailure                  // 502 (e.g., gRPC call to another service failed)
	TypeRateLimitExceeded                  // 429 (Too Many Requests)
)

// AppError là cấu trúc lỗi chuẩn của chúng ta.
type AppError struct {
	Type    ErrorType // Loại lỗi để mapping
	Code    ErrorCode // Mã lỗi để client có thể xử lý (e.g., "PRODUCT_NOT_FOUND")
	Message string    // Thông điệp cho người dùng/dev
	err     error     // Lỗi gốc (để wrapping)
}

// Implement a standard error interface.
func (e *AppError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.err)
	}
	return e.Message
}

// Unwrap để hỗ trợ error chaining (sử dụng với errors.Is, errors.As).
func (e *AppError) Unwrap() error {
	return e.err
}

// Wrap thêm một lỗi gốc vào AppError.
func (e *AppError) Wrap(err error) *AppError {
	e.err = err
	return e
}

// Các hàm khởi tạo (constructor) để dễ dàng tạo lỗi.
func New(code ErrorCode, message string, errType ErrorType) *AppError {
	return &AppError{
		Type:    errType,
		Code:    code,
		Message: message,
	}
}

func NewNotFound(resource string, id string) *AppError {
	return &AppError{
		Type:    TypeNotFound,
		Code:    CodeNotFound, // <-- Sử dụng hằng số chung
		Message: fmt.Sprintf("%s with id '%s' not found", resource, id),
	}
}

func NewValidation(message string, details error) *AppError {
	return &AppError{
		Type:    TypeValidation,
		Code:    CodeValidation, // <-- Sử dụng hằng số
		Message: message,
		err:     details,
	}
}

func NewConflict(resource string, id string) *AppError {
	return &AppError{
		Type:    TypeConflict,
		Code:    CodeConflict, // <-- Sử dụng hằng số
		Message: fmt.Sprintf("%s with id '%s' already exists", resource, id),
	}
}

func NewUnauthorized(message string) *AppError {
	return &AppError{
		Type:    TypeUnauthorized,
		Code:    CodeUnauthorized, // <-- Sử dụng hằng số
		Message: message,
	}
}

func NewForbidden(message string) *AppError {
	return &AppError{
		Type:    TypeForbidden,
		Code:    CodeForbidden, // <-- Sử dụng hằng số
		Message: message,
	}
}

func NewInternal(message string) *AppError {
	return &AppError{
		Type:    TypeInternal,
		Code:    CodeInternal, // <-- Sử dụng hằng số
		Message: message,
	}
}

func NewDependencyFailure(message string) *AppError {
	return &AppError{
		Type:    TypeDependencyFailure,
		Code:    CodeDependencyFailure, // <-- Sử dụng hằng số
		Message: message,
	}
}

func NewBadRequest(message string, details error) *AppError {
	return &AppError{
		Type:    TypeValidation,
		Code:    CodeBadRequest,
		Message: message,
		err:     details,
	}
}

// GetType là hàm helper để lấy ErrorType từ một lỗi bất kỳ.
func GetType(err error) ErrorType {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type
	}
	// Mặc định là lỗi hệ thống nếu không phải là AppError
	return TypeInternal
}

func NewRateLimitExceeded(message string) *AppError {
	return &AppError{
		Type:    TypeRateLimitExceeded,
		Code:    CodeRateLimitExceeded,
		Message: message,
	}
}

func NewTokenExpired(message string) *AppError {
	return &AppError{
		Type:    TypeUnauthorized,
		Code:    CodeTokenExpired,
		Message: message,
	}
}

func NewTokenInvalid(message string) *AppError {
	return &AppError{
		Type:    TypeUnauthorized,
		Code:    CodeTokenInvalid,
		Message: message,
	}
}
