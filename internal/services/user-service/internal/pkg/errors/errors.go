package errors

import "errors"

// Authentication and authorization errors
var (
	ErrTokenInvalid   = errors.New("token is invalid")
	ErrTokenExpired   = errors.New("token has expired")
	ErrTokenMalformed = errors.New("token is malformed")
	ErrTokenNotFound  = errors.New("token not found")
	ErrUnauthorized   = errors.New("unauthorized access")
	ErrInvalidClaims  = errors.New("invalid token claims")
)

// User errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrUserBlocked        = errors.New("user account is blocked")
)

// Validation errors
var (
	ErrInvalidInput    = errors.New("invalid input data")
	ErrMissingField    = errors.New("required field is missing")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidPassword = errors.New("invalid password format")
	ErrPasswordTooWeak = errors.New("password is too weak")
)

// Database errors
var (
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrDatabaseQuery      = errors.New("database query failed")
	ErrRecordNotFound     = errors.New("record not found")
	ErrDuplicateEntry     = errors.New("duplicate entry")
)

// Cache errors
var (
	ErrCacheConnection = errors.New("cache connection failed")
	ErrCacheNotFound   = errors.New("cache entry not found")
	ErrCacheExpired    = errors.New("cache entry expired")
)

// General errors
var (
	ErrInternalServer     = errors.New("internal server error")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrTimeout            = errors.New("request timeout")
	ErrRateLimited        = errors.New("rate limit exceeded")
)
