package constant

const (
	ContextKeyUserID    = "user_id"
	ContextKeyUserEmail = "user_email"
	ContextKeyUserRole  = "user_role"
)

type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleSeller   UserRole = "seller"
	UserRoleCustomer UserRole = "customer"
	UserRoleShipper  UserRole = "shipper"
)
