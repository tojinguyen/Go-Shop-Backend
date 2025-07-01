package constant

type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleSeller   UserRole = "seller"
	UserRoleCustomer UserRole = "customer"
	UserRoleShipper  UserRole = "shipper"
)
