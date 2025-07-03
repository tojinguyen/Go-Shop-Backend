package constant

const (
	// HTTP Status Messages
	StatusOK      = "OK"
	StatusCreated = "Created"
	StatusUpdated = "Updated"
	StatusDeleted = "Deleted"

	// Error Codes
	ErrorCodeValidation    = "VALIDATION_ERROR"
	ErrorCodeNotFound      = "NOT_FOUND"
	ErrorCodeDuplicate     = "DUPLICATE_ENTRY"
	ErrorCodeUnauthorized  = "UNAUTHORIZED"
	ErrorCodeForbidden     = "FORBIDDEN"
	ErrorCodeInternalError = "INTERNAL_ERROR"

	// Shop related constants
	ShopStatusActive    = "active"
	ShopStatusInactive  = "inactive"
	ShopStatusSuspended = "suspended"
)
